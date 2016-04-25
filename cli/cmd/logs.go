package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View output from containers",
	Run:   app.WithProject(factory, app.ProjectLog),
}

func init() {
	RootCmd.AddCommand(logsCmd)

	logsCmd.Flags().Int("lines", 100, "number of lines to tail")
	logsCmd.Flags().Bool("Follow", false, "Follow log output")

}
