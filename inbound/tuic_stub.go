//go:build !with_quic

package inbound

import (
	"context"

	"github.com/tim06/sing-box/adapter"
	C "github.com/tim06/sing-box/constant"
	"github.com/tim06/sing-box/log"
	"github.com/tim06/sing-box/option"
)

func NewTUIC(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.TUICInboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}
