package cmd

import (
	"strconv"

	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	cleanupCmd = &cobra.Command{
		Use: "cleanup",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.LogInfo.Printf("Cleaning up...\n")
			config := rootContext.Config
			cleanup := config.Cron.Cleanup
			if cleanup.Enabled {
				args := []string{"forget", "--prune"}
				if cleanup.Keep.Last > 0 {
					args = append(args, "--keep-last", strconv.Itoa(cleanup.Keep.Last))
				}
				if cleanup.Keep.Daily > 0 {
					args = append(args, "--keep-daily", strconv.Itoa(cleanup.Keep.Daily))
				}
				if cleanup.Keep.Weekly > 0 {
					args = append(args, "--keep-weekly", strconv.Itoa(cleanup.Keep.Weekly))
				}
				if cleanup.Keep.Monthly > 0 {
					args = append(args, "--keep-monthly", strconv.Itoa(cleanup.Keep.Monthly))
				}
				if cleanup.Keep.Yearly > 0 {
					args = append(args, "--keep-yearly", strconv.Itoa(cleanup.Keep.Yearly))
				}

				if err := internal.ExecRestic(rootContext, args...); err != nil {
					return err
				}
			}

			return nil
		},
	}
)
