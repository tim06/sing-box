//go:build with_wireguard

package include

import (
	"github.com/FreeVPNProxySecure/sing-box/adapter/endpoint"
	"github.com/FreeVPNProxySecure/sing-box/adapter/outbound"
	"github.com/FreeVPNProxySecure/sing-box/protocol/wireguard"
)

func registerWireGuardOutbound(registry *outbound.Registry) {
	wireguard.RegisterOutbound(registry)
}

func registerWireGuardEndpoint(registry *endpoint.Registry) {
	wireguard.RegisterEndpoint(registry)
}
