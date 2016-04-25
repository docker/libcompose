package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build or rebuild services.",
	Long: `Build or rebuild services.

Services are built once and then tagged as 'project_service',
e.g. 'composetest_db'. If you change a service's 'Dockerfile' or the
contents of its build directory, you can run 'docker-compose build' to rebuild it.
`,
	Run: app.WithProject(factory, app.ProjectBuild),
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.SetUsageTemplate(subCommandTemplate)

	//Flags
	buildCmd.Flags().BoolP("force-rm", "", false, "Always remove intermediate containers.")
	buildCmd.Flags().BoolP("no-cache", "", false, "Do not use cache when building the image.")
	buildCmd.Flags().BoolP("pull", "", false, "Always attempt to pull a newer version of the image.")

}
