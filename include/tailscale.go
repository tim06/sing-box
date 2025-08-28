//go:build with_tailscale

package include

import (
	"github.com/tim06/sing-box/adapter/endpoint"
	"github.com/tim06/sing-box/adapter/service"
	"github.com/tim06/sing-box/dns"
	"github.com/tim06/sing-box/protocol/tailscale"
	"github.com/tim06/sing-box/service/derp"
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
