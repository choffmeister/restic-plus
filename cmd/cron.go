package cmd

import (
	"io/ioutil"
	"log"

	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	cronCmd = &cobra.Command{
		Use: "cron",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.LogInfo = log.New(ioutil.Discard, "", 0)
			if err := backupCmd.RunE(cmd, []string{}); err != nil {
				return err
			}
			if err := cleanupCmd.RunE(cmd, []string{}); err != nil {
				return err
			}
			return nil
		},
	}
)
