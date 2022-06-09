package cmd

import (
	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use: "backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, target := range rootContext.Config.Targets {
				if err := internal.Restic(rootContext, "backup", target); err != nil {
					return err
				}
			}

			return nil
		},
	}
)
