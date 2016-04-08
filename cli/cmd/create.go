package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create services",
	Run:   app.WithProject(factory, app.ProjectCreate),
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.SetUsageTemplate(subCommandTemplate)

	createCmd.Flags().Bool("no-recreate", false, "If containers already exist, don't recreate them. Incompatible with --force-recreate.")
	createCmd.Flags().Bool("force-recreate", false, "Recreate containers even if their configuration and image haven't changed. Incompatible with --no-recreate.")
	createCmd.Flags().Bool("no-build", false, "Don't build an image, even if it's missing.")
}
