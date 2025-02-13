//go:build !with_quic

package inbound

import (
	C "github.com/tim06/sing-box/constant"
)

func (n *Naive) configureHTTP3Listener() error {
	return C.ErrQUICNotIncluded
}
