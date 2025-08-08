package main

import (
	"github.com/sukhera/uptime-monitor/cmd"
)

// Version information - will be set at build time via ldflags
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Set build information in the root command
	cmd.SetBuildInfo(Version, Commit, BuildDate)

	// Execute the root command
	cmd.Execute()
}
