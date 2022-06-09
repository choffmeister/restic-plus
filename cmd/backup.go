package cmd

import (
	"fmt"
	"strconv"

	"github.com/choffmeister/restic-plus/internal"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use: "backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := rootContext.Config
			targets := config.Targets
			bandwidth := config.Bandwidth
			failed := false

			for _, target := range targets {
				internal.LogInfo.Printf("Backing up %s...\n", target.Implementation.String())
				_, source, err := target.Implementation.Pre()
				defer target.Implementation.Post()
				if err != nil {
					internal.LogError.Printf("Pre for target %s failed: %v", target.Type, err)
					failed = true
					continue
				}

				args := []string{"backup", source}
				if bandwidth.Download > 0 {
					args = append(args, "--limit-downlowd", strconv.Itoa(bandwidth.Download))
				}
				if bandwidth.Upload > 0 {
					args = append(args, "--limit-upload", strconv.Itoa(bandwidth.Upload))
				}
				if err := internal.ExecRestic(rootContext, args...); err != nil {
					internal.LogError.Printf("Backup of target %s failed: %v", target.Type, err)
					failed = true
				}
			}

			if failed {
				return fmt.Errorf("some backup targets have failed")
			}
			return nil
		},
	}
)
