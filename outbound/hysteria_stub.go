//go:build !with_quic

package outbound

import (
	"context"

	"github.com/tim06/sing-box/adapter"
	C "github.com/tim06/sing-box/constant"
	"github.com/tim06/sing-box/log"
	"github.com/tim06/sing-box/option"
)

func NewHysteria(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.HysteriaOutboundOptions) (adapter.Outbound, error) {
	return nil, C.ErrQUICNotIncluded
}

func NewHysteria2(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.Hysteria2OutboundOptions) (adapter.Outbound, error) {
	return nil, C.ErrQUICNotIncluded
}
