package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers",
	Long: `Builds, (re)creates, starts, and attaches to containers for a service.

Unless they are already running, this command also starts any linked services.

The "docker-compose up" command aggregates the output of each container. When
the command exits, all containers are stopped. Running "docker-compose up -d"
starts the containers in the background and leaves them running.

If there are existing containers for a service, and the service's configuration
or image was changed after the container's creation, "docker-compose up" picks
up the changes by stopping and recreating the containers (preserving mounted
volumes). To prevent Compose from picking up changes, use the "--no-recreate"
flag.

If you want to force Compose to stop and recreate all containers, use the
"--force-recreate" flag`,
	Run: app.WithProject(factory, app.ProjectUp),
}

func init() {
	RootCmd.AddCommand(upCmd)
	upCmd.SetUsageTemplate(subCommandTemplate)

	upCmd.Flags().BoolP("detach", "d", false, "Do not block and log")
	upCmd.Flags().Bool("no-build", false, "Don't build an image, even if it's missing.")
	upCmd.Flags().Bool("no-recreate", false, "If containers already exist, don't recreate them. Incompatible with --force-recreate.")
	upCmd.Flags().Bool("force-recreate", false, "Recreate containers even if their configuration and image haven't changed. Incompatible with --no-recreate.")
}
