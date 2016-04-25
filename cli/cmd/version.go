package cmd

import (
	"fmt"
	"github.com/docker/libcompose/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the Docker-Compose version information",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("Version: ", version.VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
