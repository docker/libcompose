package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start services",
	Long:  "Start existing containers.",
	Run:   app.WithProject(factory, app.ProjectStart),
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.SetUsageTemplate(subCommandTemplate)
}
