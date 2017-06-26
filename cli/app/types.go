package app

import (
	"github.com/portainer/libcompose/project"
	"github.com/urfave/cli"
)

// ProjectFactory is an interface that helps creating libcompose project.
type ProjectFactory interface {
	// Create creates a libcompose project from the command line options (urfave cli context).
	Create(c *cli.Context) (project.APIProject, error)
}
