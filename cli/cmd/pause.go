package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause services",
	Run:   app.WithProject(factory, app.ProjectPause),
}

func init() {
	RootCmd.AddCommand(pauseCmd)
	pauseCmd.SetUsageTemplate(subCommandTemplate)
}
