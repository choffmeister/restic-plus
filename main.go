package main

import (
	"os"

	"github.com/choffmeister/restic-plus/cmd"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	cmd.Version = cmd.FullVersion{Version: version, Commit: commit, Date: date, BuiltBy: builtBy}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
