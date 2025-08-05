package main

import (
	"github.com/sukhera/uptime-monitor/cmd"
)

// Version will be set at build time via ldflags
var Version = "dev"

func main() {
	// Set version in the root command
	cmd.SetVersion(Version)

	// Execute the root command
	cmd.Execute()
}
