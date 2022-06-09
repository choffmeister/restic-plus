package cmd

import (
	"strconv"

	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	cronCmd = &cobra.Command{
		Use: "cron",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := backupCmd.RunE(cmd, []string{}); err != nil {
				return err
			}

			cleanup := rootContext.Config.Cron.Cleanup
			if cleanup.Enabled {
				forgetArgs := []string{"forget", "--prune"}
				if cleanup.Keep.Last > 0 {
					forgetArgs = append(forgetArgs, "--keep-last", strconv.Itoa(cleanup.Keep.Last))
				}
				if cleanup.Keep.Daily > 0 {
					forgetArgs = append(forgetArgs, "--keep-daily", strconv.Itoa(cleanup.Keep.Daily))
				}
				if cleanup.Keep.Weekly > 0 {
					forgetArgs = append(forgetArgs, "--keep-weekly", strconv.Itoa(cleanup.Keep.Weekly))
				}
				if cleanup.Keep.Monthly > 0 {
					forgetArgs = append(forgetArgs, "--keep-monthly", strconv.Itoa(cleanup.Keep.Monthly))
				}
				if cleanup.Keep.Yearly > 0 {
					forgetArgs = append(forgetArgs, "--keep-yearly", strconv.Itoa(cleanup.Keep.Yearly))
				}

				if err := internal.Restic(rootContext, forgetArgs...); err != nil {
					return err
				}
			}

			return nil
		},
	}
)
