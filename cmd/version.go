package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(VersionCmd)
}

var (
	// Version of command
	Version = "dev-master"

	// BuildTime in UTC
	// create by following command
	// $ TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'
	BuildTime = "undefined"

	// GitHash .
	GitHash = "undefined"
)

// VersionCmd print command version
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this binary's version, build time and git hash of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git Hash:   %s\n", GitHash)
		fmt.Printf("Build Time: %s\n", BuildTime)
	},
}
