package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart services",
	Run:   app.WithProject(factory, app.ProjectRestart),
}

func init() {
	RootCmd.AddCommand(restartCmd)
	restartCmd.SetUsageTemplate(subCommandTemplate)

	restartCmd.Flags().IntP("timeout", "t", 10, "Specify a shutdown timeout in seconds.")

}
