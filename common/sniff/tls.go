package sniff

import (
	"context"
	"crypto/tls"
	"io"

	"github.com/tim06/sing-box/adapter"
	C "github.com/tim06/sing-box/constant"
	"github.com/sagernet/sing/common/bufio"
)

func TLSClientHello(ctx context.Context, metadata *adapter.InboundContext, reader io.Reader) error {
	var clientHello *tls.ClientHelloInfo
	err := tls.Server(bufio.NewReadOnlyConn(reader), &tls.Config{
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			clientHello = argHello
			return nil, nil
		},
	}).HandshakeContext(ctx)
	if clientHello != nil {
		metadata.Protocol = C.ProtocolTLS
		metadata.Domain = clientHello.ServerName
		return nil
	}
	return err
}
