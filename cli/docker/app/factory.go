package app

import (
	"github.com/portainer/libcompose/cli/logger"
	"github.com/portainer/libcompose/docker"
	"github.com/portainer/libcompose/docker/ctx"
	"github.com/portainer/libcompose/project"
	"github.com/urfave/cli"
)

// ProjectFactory is a struct that holds the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create(c *cli.Context) (project.APIProject, error) {
	context := &ctx.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context, c)
	return docker.NewProject(context, nil)
}
