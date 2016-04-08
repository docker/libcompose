package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill containers",
	Long:  "Force stop service containers.",
	Run:   app.WithProject(factory, app.ProjectKill),
}

func init() {
	RootCmd.AddCommand(killCmd)
	killCmd.SetUsageTemplate(subCommandTemplate)

	killCmd.Flags().StringP("signal", "s", "SIGKILL", "SIGNAL to send to the container.")

}
