package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// unpauseCmd represents the unpause command
var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause services",
	Run:   app.WithProject(factory, app.ProjectUnpause),
}

func init() {
	RootCmd.AddCommand(unpauseCmd)
	unpauseCmd.SetUsageTemplate(subCommandTemplate)
}
