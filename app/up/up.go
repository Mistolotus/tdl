package up

import (
	"context"

	"github.com/Mistolotus/tdl/app/internal/tgc"
	"github.com/Mistolotus/tdl/pkg/consts"
	"github.com/Mistolotus/tdl/pkg/kv"
	"github.com/Mistolotus/tdl/pkg/uploader"
	"github.com/fatih/color"
	"github.com/gotd/td/telegram"
	"github.com/spf13/viper"
)

type Options struct {
	Chat     string
	Paths    []string
	Excludes []string
}

func WalkFile(opts *Options) ([]*file, error) {
	files, err := walk(opts.Paths, opts.Excludes)
	if err != nil {
		return nil, err
	}

	color.Blue("Files count: %d", len(files))
	return files, nil
}

func Run(ctx context.Context, opts *Options) error {
	files, err := WalkFile(opts)
	if err != nil {
		return err
	}

	c, kvd, err := tgc.NoLogin(ctx)
	if err != nil {
		return err
	}

	return ExecUp(c, kvd, files, ctx, opts)
}

func ExecUp(c *telegram.Client, kvd kv.KV, files []*file, ctx context.Context, opts *Options) error {
	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) error {
		options := &uploader.Options{
			Client:   c.API(),
			KV:       kvd,
			PartSize: viper.GetInt(consts.FlagPartSize),
			Threads:  viper.GetInt(consts.FlagThreads),
			Iter:     newIter(files),
		}
		return uploader.New(options).Upload(ctx, opts.Chat, viper.GetInt(consts.FlagLimit))
	})
}
