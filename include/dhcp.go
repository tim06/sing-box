//go:build with_dhcp

package include

import (
	"github.com/FreeVPNProxySecure/sing-box/dns"
	"github.com/FreeVPNProxySecure/sing-box/dns/transport/dhcp"
)

func registerDHCPTransport(registry *dns.TransportRegistry) {
	dhcp.RegisterTransport(registry)
}
