package docker

import (
	"github.com/Sirupsen/logrus"

	"github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
)

// ComposeVersion is name of docker-compose.yml file syntax supported version
const ComposeVersion = "1.5.0"

// NewProject creates a Project with the specified context.
func NewProject(context *Context) (project.APIProject, error) {
	if context.ResourceLookup == nil {
		context.ResourceLookup = &lookup.FileConfigLookup{}
	}

	if context.EnvironmentLookup == nil {
		context.EnvironmentLookup = &lookup.OsEnvLookup{}
	}

	if context.AuthLookup == nil {
		context.AuthLookup = &ConfigAuthLookup{context}
	}

	if context.ServiceFactory == nil {
		context.ServiceFactory = &ServiceFactory{
			context: context,
		}
	}

	if context.ClientFactory == nil {
		factory, err := project.NewDefaultClientFactory(client.Options{})
		if err != nil {
			return nil, err
		}
		context.ClientFactory = factory
	}

	p := project.NewProject(&context.Context)

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.open(); err != nil {
		logrus.Errorf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	return p, err
}
