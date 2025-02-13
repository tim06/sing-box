//go:build !with_wireguard

package outbound

import (
	"context"

	"github.com/tim06/sing-box/adapter"
	"github.com/tim06/sing-box/log"
	"github.com/tim06/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func NewWireGuard(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.WireGuardOutboundOptions) (adapter.Outbound, error) {
	return nil, E.New(`WireGuard is not included in this build, rebuild with -tags with_wireguard`)
}
