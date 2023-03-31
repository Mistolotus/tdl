package tgc

import (
	"context"
	"fmt"
	"github.com/gotd/contrib/middleware/floodwait"
	tdclock "github.com/gotd/td/clock"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/Mistolotus/tdl/pkg/clock"
	"github.com/Mistolotus/tdl/pkg/consts"
	"github.com/Mistolotus/tdl/pkg/key"
	"github.com/Mistolotus/tdl/pkg/kv"
	"github.com/Mistolotus/tdl/pkg/logger"
	"github.com/Mistolotus/tdl/pkg/storage"
	"github.com/Mistolotus/tdl/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func New(ctx context.Context, login bool, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	kvd, err := kv.New(kv.Options{
		Path: consts.KVPath,
		NS:   viper.GetString(consts.FlagNamespace),
	})
	if err != nil {
		return nil, nil, err
	}

	_clock := tdclock.System
	if ntp := viper.GetString(consts.FlagNTP); ntp != "" {
		_clock, err = clock.New()
		if err != nil {
			return nil, nil, err
		}
	}

	mode, err := kvd.Get(key.App())
	if err != nil {
		mode = []byte(consts.AppBuiltin)
	}
	app, ok := consts.Apps[string(mode)]
	if !ok {
		return nil, nil, fmt.Errorf("can't find app: %s, please try re-login", mode)
	}

	logger.From(ctx).Info("New telegram client",
		zap.Int("app", app.AppID),
		zap.String("mode", string(mode)),
		zap.Bool("is_login", login))

	return telegram.NewClient(app.AppID, app.AppHash, telegram.Options{
		Resolver: dcs.Plain(dcs.PlainOptions{
			Dial: utils.Proxy.GetDial(viper.GetString(consts.FlagProxy)).DialContext,
		}),
		Device:         consts.Device,
		SessionStorage: storage.NewSession(kvd, login),
		RetryInterval:  time.Second,
		MaxRetries:     10,
		DialTimeout:    10 * time.Second,
		Middlewares:    middlewares,
		Clock:          _clock,
		Logger:         logger.From(ctx).Named("td"),
	}), kvd, nil
}

func NoLogin(ctx context.Context, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	return New(ctx, false, append(middlewares, floodwait.NewSimpleWaiter())...)
}

func Login(ctx context.Context, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	return New(ctx, true, append(middlewares, floodwait.NewSimpleWaiter())...)
}
