package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls service images",
	Run:   app.WithProject(factory, app.ProjectPull),
}

func init() {
	RootCmd.AddCommand(pullCmd)
	pullCmd.SetUsageTemplate(subCommandTemplate)

}
