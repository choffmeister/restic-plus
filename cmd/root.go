package cmd

import (
	"io/ioutil"
	"log"

	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	rootCmd = &cobra.Command{
		Use: "restic-plus",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !verbose {
				internal.Debug = log.New(ioutil.Discard, "", log.LstdFlags)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(rawCmd)

	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(snapshotsCmd)
}
