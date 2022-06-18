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
			stanzas := config.Stanzas
			bandwidth := config.Bandwidth
			failed := false

			for _, stanza := range stanzas {
				internal.LogInfo.Printf("Backing up %s...\n", stanza.Implementation.String())
				sources, err := stanza.Implementation.Pre()
				defer stanza.Implementation.Post()

				if err != nil {
					internal.LogError.Printf("Pre for stanza %s failed: %v", stanza.Type, err)
					failed = true
					continue
				}

				for _, source := range sources {
					args := []string{"backup", source, "--json"}
					if bandwidth.Download > 0 {
						args = append(args, "--limit-downlowd", strconv.Itoa(bandwidth.Download))
					}
					if bandwidth.Upload > 0 {
						args = append(args, "--limit-upload", strconv.Itoa(bandwidth.Upload))
					}
					if err := rootContext.ExecRestic(args...); err != nil {
						internal.LogError.Printf("Backup of stanza %s failed: %v", stanza.Type, err)
						failed = true
					}
				}
			}

			if failed {
				return fmt.Errorf("some backup stanzas have failed")
			}
			return nil
		},
	}
)
