package cmd

import (
	"fmt"
	"github.com/Mistolotus/tdl/app/archive"
	"github.com/spf13/cobra"
	"time"
)

func NewBackup() *cobra.Command {
	var dst string

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup your data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dst == "" {
				dst = fmt.Sprintf("tdl-backup-%s.zip", time.Now().Format("2006-01-02-15_04_05"))
			}

			return archive.Backup(cmd.Context(), dst)
		},
	}

	cmd.Flags().StringVarP(&dst, "dst", "d", "", "destination file path. Default: tdl-backup-<time>.zip")

	return cmd
}

func NewRecover() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "recover",
		Short: "Recover your data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return archive.Recover(cmd.Context(), file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "backup file path")

	return cmd
}
