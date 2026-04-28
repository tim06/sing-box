//go:build !android

package tailscale

import "github.com/FreeVPNProxySecure/sing-box/experimental/libbox/platform"

func setAndroidProtectFunc(platformInterface platform.Interface) {
}
