package libbox

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"
	"sync"
	"syscall"
	"time"

	tun "github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/control"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/service"
	"github.com/sagernet/sing/service/filemanager"
	"github.com/sagernet/sing/service/pause"
	box "github.com/tim06/sing-box"
	"github.com/tim06/sing-box/adapter"
	"github.com/tim06/sing-box/common/process"
	"github.com/tim06/sing-box/common/urltest"
	C "github.com/tim06/sing-box/constant"
	"github.com/tim06/sing-box/experimental/deprecated"
	"github.com/tim06/sing-box/experimental/libbox/internal/procfs"
	"github.com/tim06/sing-box/experimental/libbox/platform"
	"github.com/tim06/sing-box/log"
	"github.com/tim06/sing-box/option"
	"github.com/tim06/sing-box/protocol/group"
)

type BoxService struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	urlTestHistoryStorage adapter.URLTestHistoryStorage
	instance              *box.Box
	clashServer           adapter.ClashServer
	pauseManager          pause.Manager

	iOSPauseFields
}

func NewService(configContent string, platformInterface PlatformInterface) (*BoxService, error) {
	ctx := BaseContext(platformInterface)
	ctx = filemanager.WithDefault(ctx, sWorkingPath, sTempPath, sUserID, sGroupID)
	service.MustRegister[deprecated.Manager](ctx, new(deprecatedManager))
	options, err := parseConfig(ctx, configContent)
	if err != nil {
		return nil, err
	}
	runtimeDebug.FreeOSMemory()
	ctx, cancel := context.WithCancel(ctx)
	urlTestHistoryStorage := urltest.NewHistoryStorage()
	ctx = service.ContextWithPtr(ctx, urlTestHistoryStorage)
	platformWrapper := &platformInterfaceWrapper{
		iif:       platformInterface,
		useProcFS: platformInterface.UseProcFS(),
	}
	service.MustRegister[platform.Interface](ctx, platformWrapper)
	instance, err := box.New(box.Options{
		Context:           ctx,
		Options:           options,
		PlatformLogWriter: platformWrapper,
	})
	if err != nil {
		cancel()
		return nil, E.Cause(err, "create service")
	}
	runtimeDebug.FreeOSMemory()
	return &BoxService{
		ctx:                   ctx,
		cancel:                cancel,
		instance:              instance,
		urlTestHistoryStorage: urlTestHistoryStorage,
		pauseManager:          service.FromContext[pause.Manager](ctx),
		clashServer:           service.FromContext[adapter.ClashServer](ctx),
	}, nil
}

func (s *BoxService) Start() error {
	var err error
	if sFixAndroidStack {
		done := make(chan struct{})
		go func() {
			err = s.instance.Start()
			close(done)
		}()
		<-done
	} else {
		err = s.instance.Start()
	}

	if err == nil {
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			defer ticker.Stop()

			if selectErr := s.selectFirstAvailableOutbound(); selectErr != nil {
				log.Warn("Failed to select first available outbound: ", selectErr)
			}

			for {
				select {
				case <-ticker.C:
					if selectErr := s.selectFirstAvailableOutbound(); selectErr != nil {
						log.Warn("Failed to select first available outbound: ", selectErr)
					}
				case <-s.ctx.Done():
					return
				}
			}
		}()
	}

	return err
}

func (s *BoxService) Close() error {
	s.cancel()
	s.urlTestHistoryStorage.Close()
	var err error
	done := make(chan struct{})
	go func() {
		err = s.instance.Close()
		close(done)
	}()
	select {
	case <-done:
		return err
	case <-time.After(C.FatalStopTimeout):
		os.Exit(1)
		return nil
	}
}

func (s *BoxService) NeedWIFIState() bool {
	return s.instance.Router().NeedWIFIState()
}

var (
	_ platform.Interface = (*platformInterfaceWrapper)(nil)
	_ log.PlatformWriter = (*platformInterfaceWrapper)(nil)
)

type platformInterfaceWrapper struct {
	iif                    PlatformInterface
	useProcFS              bool
	networkManager         adapter.NetworkManager
	myTunName              string
	defaultInterfaceAccess sync.Mutex
	defaultInterface       *control.Interface
	isExpensive            bool
	isConstrained          bool
}

func (w *platformInterfaceWrapper) Initialize(networkManager adapter.NetworkManager) error {
	w.networkManager = networkManager
	return nil
}

func (w *platformInterfaceWrapper) UsePlatformAutoDetectInterfaceControl() bool {
	return w.iif.UsePlatformAutoDetectInterfaceControl()
}

func (w *platformInterfaceWrapper) AutoDetectInterfaceControl(fd int) error {
	return w.iif.AutoDetectInterfaceControl(int32(fd))
}

func (w *platformInterfaceWrapper) OpenTun(options *tun.Options, platformOptions option.TunPlatformOptions) (tun.Tun, error) {
	if len(options.IncludeUID) > 0 || len(options.ExcludeUID) > 0 {
		return nil, E.New("platform: unsupported uid options")
	}
	if len(options.IncludeAndroidUser) > 0 {
		return nil, E.New("platform: unsupported android_user option")
	}
	routeRanges, err := options.BuildAutoRouteRanges(true)
	if err != nil {
		return nil, err
	}
	tunFd, err := w.iif.OpenTun(&tunOptions{options, routeRanges, platformOptions})
	if err != nil {
		return nil, err
	}
	options.Name, err = getTunnelName(tunFd)
	if err != nil {
		return nil, E.Cause(err, "query tun name")
	}
	dupFd, err := dup(int(tunFd))
	if err != nil {
		return nil, E.Cause(err, "dup tun file descriptor")
	}
	options.FileDescriptor = dupFd
	w.myTunName = options.Name
	return tun.New(*options)
}

func (w *platformInterfaceWrapper) CreateDefaultInterfaceMonitor(logger logger.Logger) tun.DefaultInterfaceMonitor {
	return &platformDefaultInterfaceMonitor{
		platformInterfaceWrapper: w,
		logger:                   logger,
	}
}

func (w *platformInterfaceWrapper) Interfaces() ([]adapter.NetworkInterface, error) {
	interfaceIterator, err := w.iif.GetInterfaces()
	if err != nil {
		return nil, err
	}
	var interfaces []adapter.NetworkInterface
	for _, netInterface := range iteratorToArray[*NetworkInterface](interfaceIterator) {
		if netInterface.Name == w.myTunName {
			continue
		}
		w.defaultInterfaceAccess.Lock()
		// (GOOS=windows) SA4006: this value of `isDefault` is never used
		// Why not used?
		//nolint:staticcheck
		isDefault := w.defaultInterface != nil && int(netInterface.Index) == w.defaultInterface.Index
		w.defaultInterfaceAccess.Unlock()
		interfaces = append(interfaces, adapter.NetworkInterface{
			Interface: control.Interface{
				Index:     int(netInterface.Index),
				MTU:       int(netInterface.MTU),
				Name:      netInterface.Name,
				Addresses: common.Map(iteratorToArray[string](netInterface.Addresses), netip.MustParsePrefix),
				Flags:     linkFlags(uint32(netInterface.Flags)),
			},
			Type:        C.InterfaceType(netInterface.Type),
			DNSServers:  iteratorToArray[string](netInterface.DNSServer),
			Expensive:   netInterface.Metered || isDefault && w.isExpensive,
			Constrained: isDefault && w.isConstrained,
		})
	}
	return interfaces, nil
}

func (w *platformInterfaceWrapper) UnderNetworkExtension() bool {
	return w.iif.UnderNetworkExtension()
}

func (w *platformInterfaceWrapper) IncludeAllNetworks() bool {
	return w.iif.IncludeAllNetworks()
}

func (w *platformInterfaceWrapper) ClearDNSCache() {
	w.iif.ClearDNSCache()
}

func (w *platformInterfaceWrapper) ReadWIFIState() adapter.WIFIState {
	wifiState := w.iif.ReadWIFIState()
	if wifiState == nil {
		return adapter.WIFIState{}
	}
	return (adapter.WIFIState)(*wifiState)
}

func (w *platformInterfaceWrapper) SystemCertificates() []string {
	return iteratorToArray[string](w.iif.SystemCertificates())
}

func (w *platformInterfaceWrapper) FindProcessInfo(ctx context.Context, network string, source netip.AddrPort, destination netip.AddrPort) (*process.Info, error) {
	var uid int32
	if w.useProcFS {
		uid = procfs.ResolveSocketByProcSearch(network, source, destination)
		if uid == -1 {
			return nil, E.New("procfs: not found")
		}
	} else {
		var ipProtocol int32
		switch N.NetworkName(network) {
		case N.NetworkTCP:
			ipProtocol = syscall.IPPROTO_TCP
		case N.NetworkUDP:
			ipProtocol = syscall.IPPROTO_UDP
		default:
			return nil, E.New("unknown network: ", network)
		}
		var err error
		uid, err = w.iif.FindConnectionOwner(ipProtocol, source.Addr().String(), int32(source.Port()), destination.Addr().String(), int32(destination.Port()))
		if err != nil {
			return nil, err
		}
	}
	packageName, _ := w.iif.PackageNameByUid(uid)
	return &process.Info{UserId: uid, PackageName: packageName}, nil
}

func (w *platformInterfaceWrapper) DisableColors() bool {
	return runtime.GOOS != "android"
}

func (w *platformInterfaceWrapper) WriteMessage(level log.Level, message string) {
	w.iif.WriteLog(message)
}

func (w *platformInterfaceWrapper) SendNotification(notification *platform.Notification) error {
	return w.iif.SendNotification((*Notification)(notification))
}

func (s *BoxService) selectFirstAvailableOutbound() error {
	var selectorOutbound *group.Selector

	for _, outbound := range s.instance.Outbound().Outbounds() {
		if selector, ok := outbound.(*group.Selector); ok {
			selectorOutbound = selector
			break
		}
	}

	if selectorOutbound == nil {
		return fmt.Errorf("no selector outbound found")
	}

	for _, outboundTag := range selectorOutbound.All() {
		outbound, found := s.instance.Outbound().Outbound(outboundTag)
		if !found {
			continue
		}

		if outbound.Type() == C.TypeURLTest {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		delay, err := urltest.URLTest(ctx, "https://www.gstatic.com/generate_204", outbound)
		cancel()

		if err == nil && delay > 0 {
			if selectorOutbound.SelectOutbound(outboundTag) {
				return nil
			}
		}
	}

	return fmt.Errorf("no available outbound found")
}
