package cmd

import (
	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use: "backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := internal.Config{}
			if err := config.LoadFromFile(""); err != nil {
				return err
			}

			for _, target := range config.Targets {
				if err := internal.Restic("backup", target); err != nil {
					return err
				}
			}

			return nil
		},
	}
)
