package app

import (
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/cli/logger"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

// ProjectFactory is a struct that holds the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create(c *cli.Context) (project.APIProject, error) {
	context := &docker.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	context.Logger = logrus.New()
	Populate(context, c)

	context.ComposeFiles = c.GlobalStringSlice("file")

	if len(context.ComposeFiles) == 0 {
		context.ComposeFiles = []string{"docker-compose.yml"}
		if _, err := os.Stat("docker-compose.override.yml"); err == nil {
			context.ComposeFiles = append(context.ComposeFiles, "docker-compose.override.yml")
		}
	}

	context.ProjectName = normalizeName(c.GlobalString("project-name"))

	return docker.NewProject(context)
}

func normalizeName(name string) string {
	r := regexp.MustCompile("[^a-z0-9]+")
	return r.ReplaceAllString(strings.ToLower(name), "")
}
