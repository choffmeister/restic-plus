package cmd

import (
	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	rawCmd = &cobra.Command{
		Use: "raw",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.Restic(rootContext, args...); err != nil {
				return err
			}

			return nil
		},
	}
)
