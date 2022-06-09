package cmd

import (
	"io/ioutil"
	"log"

	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	rootCmdVerbose bool
	rootCmdConfig  string
	rootCmd        = &cobra.Command{
		Use: "restic-plus",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !rootCmdVerbose {
				internal.Debug = log.New(ioutil.Discard, "", log.LstdFlags)
			}
			context, err := internal.NewContext(rootCmdConfig)
			if err != nil {
				return err
			}
			rootContext = context
			return nil
		},
	}
	rootContext *internal.Context
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootCmdVerbose, "verbose", "v", false, "")
	rootCmd.PersistentFlags().StringVarP(&rootCmdConfig, "config", "c", "restic-plus.yaml", "")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(rawCmd)

	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(snapshotsCmd)
}
