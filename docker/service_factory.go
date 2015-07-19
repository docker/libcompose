package docker

import "github.com/docker/libcompose/project"

type ServiceFactory struct {
	context *Context
}

func (s *ServiceFactory) Create(project *project.Project, name string, serviceConfig *project.ServiceConfig) (project.Service, error) {
	return &Service{
		name:          name,
		serviceConfig: serviceConfig,
		context:       s.context,
	}, nil
}
