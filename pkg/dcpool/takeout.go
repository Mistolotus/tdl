package dcpool

import (
	"github.com/gotd/td/tg"
	"github.com/Mistolotus/tdl/pkg/takeout"
)

func (p *pool) Takeout(dc int) *tg.Client {
	return tg.NewClient(chainMiddlewares(p.invoker(dc), takeout.Middleware(p.takeout)))
}
