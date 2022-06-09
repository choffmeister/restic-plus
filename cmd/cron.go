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
			config := internal.Config{}
			if err := config.LoadFromFile(""); err != nil {
				return err
			}

			if err := backupCmd.RunE(cmd, []string{}); err != nil {
				return err
			}

			if config.Cron.Cleanup.Enabled {
				forgetArgs := []string{"forget", "--prune"}
				if config.Cron.Cleanup.Keep.Last > 0 {
					forgetArgs = append(forgetArgs, "--keep-last", strconv.Itoa(config.Cron.Cleanup.Keep.Last))
				}
				if config.Cron.Cleanup.Keep.Daily > 0 {
					forgetArgs = append(forgetArgs, "--keep-daily", strconv.Itoa(config.Cron.Cleanup.Keep.Daily))
				}
				if config.Cron.Cleanup.Keep.Weekly > 0 {
					forgetArgs = append(forgetArgs, "--keep-weekly", strconv.Itoa(config.Cron.Cleanup.Keep.Weekly))
				}
				if config.Cron.Cleanup.Keep.Monthly > 0 {
					forgetArgs = append(forgetArgs, "--keep-monthly", strconv.Itoa(config.Cron.Cleanup.Keep.Monthly))
				}
				if config.Cron.Cleanup.Keep.Yearly > 0 {
					forgetArgs = append(forgetArgs, "--keep-yearly", strconv.Itoa(config.Cron.Cleanup.Keep.Yearly))
				}

				if err := internal.Restic(forgetArgs...); err != nil {
					return err
				}
			}

			return nil
		},
	}
)
