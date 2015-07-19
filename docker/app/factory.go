package app

import (
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/command"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/logger"
	"github.com/docker/libcompose/project"
)

type ProjectFactory struct {
}

func (p *ProjectFactory) Create(c *cli.Context) (*project.Project, error) {
	context := &docker.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context, c)
	command.Populate(&context.Context, c)

	return docker.NewProject(context)
}
