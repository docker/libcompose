package cmd

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a one-off command",
	Long: `Run a one-off command on a service.

For example:

    $ docker-compose run web python manage.py shell

By default, linked services will be started, unless they are already
running. If you do not want to start linked services, use
"docker-compose run --no-deps SERVICE COMMAND [ARGS...]".`,

	Run: app.WithProject(factory, app.ProjectRun),
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.SetUsageTemplate(subCommandTemplate)

}
