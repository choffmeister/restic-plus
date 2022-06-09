package cmd

import (
	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	snapshotsCmd = &cobra.Command{
		Use: "snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.Restic("snapshots"); err != nil {
				return err
			}

			return nil
		},
	}
)
