package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Set number of containers for a service",
	Long: `
Set number of containers to run for a service.

Numbers are specified in the form "service=num" as arguments.
For example:

    $ docker-compose scale web=2 worker=3 `,
	Run: app.WithProject(factory, app.ProjectScale),
}

func init() {
	RootCmd.AddCommand(scaleCmd)
	scaleCmd.SetUsageTemplate(subCommandTemplate)

	scaleCmd.Flags().IntP("timeout", "t", 10, "Specify a shutdown timeout in seconds.")

}
