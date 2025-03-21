//go:build !android

package tailscale

import "github.com/tim06/sing-box/experimental/libbox/platform"

func setAndroidProtectFunc(platformInterface platform.Interface) {
}
