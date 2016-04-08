package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// portCmd represents the port command
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Print the public port for a port binding",
	Run:   app.WithProject(factory, app.ProjectPort),
}

func init() {
	RootCmd.AddCommand(portCmd)
	portCmd.SetUsageTemplate(subCommandTemplate)

	portCmd.Flags().String("protocol", "tcp", "tcp or udp")
	portCmd.Flags().Int("index", 1, "index of the container if there are multiple instances of a service")

}
