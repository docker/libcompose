package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Run:   app.WithProject(factory, app.ProjectPs),
}

func init() {
	RootCmd.AddCommand(psCmd)
	psCmd.SetUsageTemplate(subCommandTemplate)

	psCmd.Flags().BoolP("quite", "q", false, "Only display IDs")
}
