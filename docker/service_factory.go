package docker

import "github.com/docker/libcompose/project"

type ServiceFactory struct {
	Context *Context
}

func (s *ServiceFactory) Create(project *project.Project, name string, serviceConfig *project.ServiceConfig) (project.Service, error) {
	return NewService(name, serviceConfig, s.Context), nil
}
