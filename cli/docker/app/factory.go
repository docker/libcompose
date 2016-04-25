package app

import (
	"github.com/docker/libcompose/cli/app"
	"github.com/docker/libcompose/cli/logger"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/spf13/cobra"
)

// ProjectFactory is a struct that holds the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create(c *cobra.Command) (*project.Project, error) {
	context := &docker.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context, c)
	app.Populate(&context.Context, c)

	return docker.NewProject(context)
}
