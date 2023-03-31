package dl

import (
	"context"
	"github.com/gotd/td/telegram/peers"
	"github.com/Mistolotus/tdl/app/internal/dliter"
	"github.com/Mistolotus/tdl/pkg/dcpool"
	"github.com/Mistolotus/tdl/pkg/kv"
	"github.com/Mistolotus/tdl/pkg/logger"
	"github.com/Mistolotus/tdl/pkg/storage"
	"github.com/Mistolotus/tdl/pkg/utils"
	"go.uber.org/zap"
)

func parseURLs(ctx context.Context, pool dcpool.Pool, kvd kv.KV, urls []string) ([]*dliter.Dialog, error) {
	manager := peers.Options{Storage: storage.NewPeers(kvd)}.
		Build(pool.Client(pool.Default()))
	msgMap := make(map[int64]*dliter.Dialog)

	for _, u := range urls {
		ch, msgid, err := utils.Telegram.ParseMessageLink(ctx, manager, u)
		if err != nil {
			return nil, err
		}
		logger.From(ctx).Debug("Parse URL",
			zap.String("url", u),
			zap.Int64("peer_id", ch.ID()),
			zap.String("peer_name", ch.VisibleName()),
			zap.Int("msg", msgid))

		// init map value
		if _, ok := msgMap[ch.ID()]; !ok {
			msgMap[ch.ID()] = &dliter.Dialog{Peer: ch.InputPeer(), Messages: []int{}}
		}

		msgMap[ch.ID()].Messages = append(msgMap[ch.ID()].Messages, msgid)
	}

	// cap is at least len of map
	msgs := make([]*dliter.Dialog, 0, len(msgMap))
	for _, m := range msgMap {
		msgs = append(msgs, m)
	}

	return msgs, nil
}
