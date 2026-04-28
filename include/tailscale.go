//go:build with_tailscale

package include

import (
	"github.com/FreeVPNProxySecure/sing-box/adapter/endpoint"
	"github.com/FreeVPNProxySecure/sing-box/adapter/service"
	"github.com/FreeVPNProxySecure/sing-box/dns"
	"github.com/FreeVPNProxySecure/sing-box/protocol/tailscale"
	"github.com/FreeVPNProxySecure/sing-box/service/derp"
)

func registerTailscaleEndpoint(registry *endpoint.Registry) {
	tailscale.RegisterEndpoint(registry)
}

func registerTailscaleTransport(registry *dns.TransportRegistry) {
	tailscale.RegistryTransport(registry)
}

func registerDERPService(registry *service.Registry) {
	derp.Register(registry)
}
