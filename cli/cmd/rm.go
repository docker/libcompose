package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove stopped containers",
	Long: `Remove stopped service containers.

By default, volumes attached to containers will not be removed. You can see all
volumes with "docker volume ls".

Any data which is not in a volume will be lost.`,
	Run: app.WithProject(factory, app.ProjectDelete),
}

func init() {
	RootCmd.AddCommand(rmCmd)
	rmCmd.SetUsageTemplate(subCommandTemplate)
	rmCmd.Flags().BoolP("force", "f", false, "Don't ask to confirm removal")
	rmCmd.Flags().BoolP("volume", "v", false, "Remove volumes associated with containers")
}
