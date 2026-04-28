//go:build with_quic

package include

import (
	"github.com/FreeVPNProxySecure/sing-box/adapter/inbound"
	"github.com/FreeVPNProxySecure/sing-box/adapter/outbound"
	"github.com/FreeVPNProxySecure/sing-box/dns"
	"github.com/FreeVPNProxySecure/sing-box/dns/transport/quic"
	"github.com/FreeVPNProxySecure/sing-box/protocol/hysteria"
	"github.com/FreeVPNProxySecure/sing-box/protocol/hysteria2"
	_ "github.com/FreeVPNProxySecure/sing-box/protocol/naive/quic"
	"github.com/FreeVPNProxySecure/sing-box/protocol/tuic"
	_ "github.com/FreeVPNProxySecure/sing-box/transport/v2rayquic"
)

func registerQUICInbounds(registry *inbound.Registry) {
	hysteria.RegisterInbound(registry)
	tuic.RegisterInbound(registry)
	hysteria2.RegisterInbound(registry)
}

func registerQUICOutbounds(registry *outbound.Registry) {
	hysteria.RegisterOutbound(registry)
	tuic.RegisterOutbound(registry)
	hysteria2.RegisterOutbound(registry)
}

func registerQUICTransports(registry *dns.TransportRegistry) {
	quic.RegisterTransport(registry)
	quic.RegisterHTTP3Transport(registry)
}
