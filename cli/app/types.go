package app

import (
	"github.com/spf13/cobra"
	"github.com/docker/libcompose/project"
)

// ProjectFactory is an interface that helps creating libcompose project.
type ProjectFactory interface {
	// Create creates a libcompose project from the command line options (codegangsta cli context).
	Create(c *cobra.Command) (*project.Project, error)
}
