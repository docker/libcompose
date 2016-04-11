package app

import (
	"os"

	"github.com/spf13/cobra"
	//	"github.com/docker/libcompose/newcli/app"
	"github.com/docker/libcompose/project"
)

// Populate updates the specified project context based on command line arguments and subcommands.
func Populate(context *project.Context, c *cobra.Command) {

	args := c.Flags()

	context.ComposeFiles, _ = args.GetStringSlice("file")

	if len(context.ComposeFiles) == 0 {
		context.ComposeFiles = []string{"docker-compose.yml"}
		if _, err := os.Stat("docker-compose.override.yml"); err == nil {
			context.ComposeFiles = append(context.ComposeFiles, "docker-compose.override.yml")
		}
	}

	context.ProjectName, _ = args.GetString("project-name")

	if c.Name() == "logs" {
		context.Log = true
	} else if c.Name() == "up" || c.Name() == "create" {
		context.Log, _ = args.GetBool("d")
		context.NoRecreate, _ = args.GetBool("no-recreate")
		context.ForceRecreate, _ = args.GetBool("force-recreate")
		context.NoBuild, _ = args.GetBool("no-build")
		//	} else if c.Name() == "stop" || c.Name() == "restart" || c.Name() == "scale" {
		//		context.Timeout, _ = args.GetInt("timeout")
		//} else if c.Name() == "kill" {
		//	context.Signal, _ = args.GetBool("signal")
	} else if c.Name() == "rm" {
		context.Volume, _ = args.GetBool("v")
	} else if c.Name() == "build" {
		context.NoCache, _ = args.GetBool("no-cache")
	}
}
