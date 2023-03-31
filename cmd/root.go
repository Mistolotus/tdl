package cmd

import (
	"github.com/fatih/color"
	"github.com/Mistolotus/tdl/pkg/consts"
	"github.com/Mistolotus/tdl/pkg/logger"
	"github.com/Mistolotus/tdl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"path/filepath"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "tdl",
		Short:             "Telegram Downloader, but more than a downloader",
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// init logger
			debug, level := viper.GetBool(consts.FlagDebug), zap.InfoLevel
			if debug {
				level = zap.DebugLevel
			}
			cmd.SetContext(logger.With(cmd.Context(), logger.New(level)))

			ns := viper.GetString(consts.FlagNamespace)
			if ns != "" {
				color.Cyan("Namespace: %s", ns)
				logger.From(cmd.Context()).Info("Namespace",
					zap.String("namespace", ns))
			}
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return logger.From(cmd.Context()).Sync()
		},
	}

	cmd.AddCommand(NewVersion(), NewLogin(), NewDownload(),
		NewChat(), NewUpload(), NewBackup(), NewRecover())

	cmd.PersistentFlags().String(consts.FlagProxy, "", "proxy address, only socks5 is supported, format: protocol://username:password@host:port")
	cmd.PersistentFlags().StringP(consts.FlagNamespace, "n", "", "namespace for Telegram session")
	cmd.PersistentFlags().Bool(consts.FlagDebug, false, "enable debug mode")

	// The default parameters are consistent with the official client to reduce the probability of blocking
	// https://github.com/Mistolotus/tdl/issues/30
	cmd.PersistentFlags().IntP(consts.FlagPartSize, "s", 128*1024, "part size for transfer, max is 512*1024")
	cmd.PersistentFlags().IntP(consts.FlagThreads, "t", 4, "max threads for transfer one item")
	cmd.PersistentFlags().IntP(consts.FlagLimit, "l", 2, "max number of concurrent tasks")

	cmd.PersistentFlags().String(consts.FlagNTP, "", "ntp server host, if not set, use system time")

	_ = viper.BindPFlags(cmd.PersistentFlags())

	viper.SetEnvPrefix("tdl")
	viper.AutomaticEnv()

	generateCommandDocs(cmd)

	return cmd
}

func generateCommandDocs(cmd *cobra.Command) {
	docs := filepath.Join(consts.DocsPath, "command")
	if utils.FS.PathExists(docs) {
		if err := doc.GenMarkdownTree(cmd, docs); err != nil {
			color.Red("generate cmd docs failed: %v", err)
		}
	}
}
