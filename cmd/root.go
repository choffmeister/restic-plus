package cmd

import (
	"io/ioutil"
	"log"
	"os"

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
				internal.LogDebug = log.New(ioutil.Discard, "", 0)
				internal.LogRestic = log.New(ioutil.Discard, "", 0)
			}
			context, err := internal.NewContext(rootCmdConfig)
			if err != nil {
				return err
			}
			rootContext = context
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Usage()
			}
			internal.LogRestic = log.New(os.Stdout, "", 0)
			if err := rootContext.ExecRestic(args...); err != nil {
				return err
			}
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

	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(cronCmd)
}
