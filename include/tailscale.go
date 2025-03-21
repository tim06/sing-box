//go:build with_tailscale

package include

import (
	"github.com/tim06/sing-box/adapter/endpoint"
	"github.com/tim06/sing-box/dns"
	"github.com/tim06/sing-box/protocol/tailscale"
)

func registerTailscaleEndpoint(registry *endpoint.Registry) {
	tailscale.RegisterEndpoint(registry)
}

func registerTailscaleTransport(registry *dns.TransportRegistry) {
	tailscale.RegistryTransport(registry)
}
