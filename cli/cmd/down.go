package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers, networks, images, and volumes",
	Run:   app.WithProject(factory, app.ProjectDown),
}

func init() {
	RootCmd.AddCommand(downCmd)
	downCmd.SetUsageTemplate(subCommandTemplate)

	downCmd.Flags().IntP("timeout", "t", 10, "Specify a shutdown timeout in seconds.")
}
