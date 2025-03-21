//go:build !with_clash_api

package include

import (
	"context"

	"github.com/tim06/sing-box/adapter"
	"github.com/tim06/sing-box/experimental"
	"github.com/tim06/sing-box/log"
	"github.com/tim06/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func init() {
	experimental.RegisterClashServerConstructor(func(ctx context.Context, logFactory log.ObservableFactory, options option.ClashAPIOptions) (adapter.ClashServer, error) {
		return nil, E.New(`clash api is not included in this build, rebuild with -tags with_clash_api`)
	})
}
