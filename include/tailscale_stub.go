//go:build !with_tailscale

package include

import (
	"context"

	"github.com/FreeVPNProxySecure/sing-box/adapter"
	"github.com/FreeVPNProxySecure/sing-box/adapter/endpoint"
	"github.com/FreeVPNProxySecure/sing-box/adapter/service"
	C "github.com/FreeVPNProxySecure/sing-box/constant"
	"github.com/FreeVPNProxySecure/sing-box/dns"
	"github.com/FreeVPNProxySecure/sing-box/log"
	"github.com/FreeVPNProxySecure/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func registerTailscaleEndpoint(registry *endpoint.Registry) {
	endpoint.Register[option.TailscaleEndpointOptions](registry, C.TypeTailscale, func(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.TailscaleEndpointOptions) (adapter.Endpoint, error) {
		return nil, E.New(`Tailscale is not included in this build, rebuild with -tags with_tailscale`)
	})
}

func registerTailscaleTransport(registry *dns.TransportRegistry) {
	dns.RegisterTransport[option.TailscaleDNSServerOptions](registry, C.DNSTypeTailscale, func(ctx context.Context, logger log.ContextLogger, tag string, options option.TailscaleDNSServerOptions) (adapter.DNSTransport, error) {
		return nil, E.New(`Tailscale is not included in this build, rebuild with -tags with_tailscale`)
	})
}

func registerDERPService(registry *service.Registry) {
	service.Register[option.DERPServiceOptions](registry, C.TypeDERP, func(ctx context.Context, logger log.ContextLogger, tag string, options option.DERPServiceOptions) (adapter.Service, error) {
		return nil, E.New(`DERP is not included in this build, rebuild with -tags with_tailscale`)
	})
}
