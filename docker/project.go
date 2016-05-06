package docker

import (
	"github.com/docker/libcompose/logger"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
)

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

	var log logger.Logger = &logger.DefaultLogger{}
	if context.Logger != nil {
		log = context.Logger
	} else {
		context.Logger = log
	}

	if context.ClientFactory == nil {
		factory, err := NewDefaultClientFactory(ClientOpts{})
		if err != nil {
			return nil, err
		}
		context.ClientFactory = factory
	}

	p := project.NewProject(&context.Context, log)

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.open(); err != nil {
		log.Errorf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	return p, err
}
