package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop services",
	Long: `Stop running containers without removing them.

They can be started again with "docker-compose start".`,

	Run: app.WithProject(factory, app.ProjectDown),
}

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.SetUsageTemplate(subCommandTemplate)

	stopCmd.Flags().IntP("timeout", "t", 10, "Specify a shutdown timeout in seconds.")

}
