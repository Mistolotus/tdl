package cmd

import (
	"github.com/Mistolotus/tdl/app/up"
	"github.com/Mistolotus/tdl/pkg/consts"
	"github.com/Mistolotus/tdl/pkg/logger"
	"github.com/spf13/cobra"
)

func NewUpload() *cobra.Command {
	var opts up.Options

	cmd := &cobra.Command{
		Use:     "upload",
		Aliases: []string{"up"},
		Short:   "Upload anything to Telegram",
		RunE: func(cmd *cobra.Command, args []string) error {
			return up.Run(logger.Named(cmd.Context(), "up"), &opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Chat, "chat", "c", "", "chat id or domain, and empty means 'Saved Messages'")
	cmd.Flags().StringSliceVarP(&opts.Paths, consts.FlagUpPath, "p", []string{}, "dirs or files")
	cmd.Flags().StringSliceVarP(&opts.Excludes, consts.FlagUpExcludes, "e", []string{}, "exclude the specified file extensions")

	return cmd
}
