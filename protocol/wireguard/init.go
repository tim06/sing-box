package wireguard

import (
	"github.com/tim06/sing-box/common/dialer"
	"github.com/sagernet/wireguard-go/conn"
)

func init() {
	dialer.WgControlFns = conn.ControlFns
}
