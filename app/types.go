package app

import (
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/project"
)

type ProjectFactory interface {
	Create(c *cli.Context) (*project.Project, error)
}
