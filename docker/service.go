package docker

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/utils"
)

// Service is a project.Service implementations.
type Service struct {
	name          string
	serviceConfig *project.ServiceConfig
	context       *Context
	imageName     string
}

// Name returns the service name.
func (s *Service) Name() string {
	return s.name
}

// Config returns the configuration of the service (project.ServiceConfig).
func (s *Service) Config() *project.ServiceConfig {
	return s.serviceConfig
}

// DependentServices returns the dependent services (as an array of ServiceRelationship) of the service.
func (s *Service) DependentServices() []project.ServiceRelationship {
	return project.DefaultDependentServices(s.context.Project, s)
}

// Create implements Service.Create.
func (s *Service) Create() error {
	_, err := s.createOne()
	return err
}

func (s *Service) collectContainers() ([]*Container, error) {
	client := s.context.ClientFactory.Create(s)
	containers, err := GetContainersByFilter(client, SERVICE.Eq(s.name), PROJECT.Eq(s.context.Project.Name))
	if err != nil {
		return nil, err
	}

	result := []*Container{}

	for _, container := range containers {
		result = append(result, NewContainer(client, container.Labels[NAME.Str()], s))
	}

	return result, nil
}

func (s *Service) createOne() (*Container, error) {
	containers, err := s.constructContainers(true, 1)
	if err != nil {
		return nil, err
	}

	return containers[0], err
}

// Build implements Service.Build. If an imageName is specified or if the context has
// no build to work with it will do nothing. Otherwise it will try to build
// the image and returns an error if any.
func (s *Service) Build() error {
	_, err := s.build()
	return err
}

func (s *Service) build() (string, error) {
	if s.imageName != "" {
		return s.imageName, nil
	}

	if s.context.Builder == nil {
		s.imageName = s.Config().Image
	} else {
		var err error
		s.imageName, err = s.context.Builder.Build(s.context.Project, s)
		if err != nil {
			return "", err
		}
	}

	return s.imageName, nil
}

func (s *Service) constructContainers(create bool, count int) ([]*Container, error) {
	result, err := s.collectContainers()
	if err != nil {
		return nil, err
	}

	client := s.context.ClientFactory.Create(s)

	var namer Namer

	if s.serviceConfig.ContainerName != "" {
		if count > 1 {
			logrus.Warnf(`The "%s" service is using the custom container name "%s". Docker requires each container to have a unique name. Remove the custom name to scale the service.`, s.name, s.serviceConfig.ContainerName)
		}
		namer = NewSingleNamer(s.serviceConfig.ContainerName)
	} else {
		namer = NewNamer(client, s.context.Project.Name, s.name)
	}

	defer namer.Close()

	for i := len(result); i < count; i++ {
		containerName := namer.Next()

		c := NewContainer(client, containerName, s)

		if create {
			imageName, err := s.build()
			if err != nil {
				return nil, err
			}

			dockerContainer, err := c.Create(imageName)
			if err != nil {
				return nil, err
			}
			logrus.Debugf("Created container %s: %v", dockerContainer.Id, dockerContainer.Names)
		}

		result = append(result, c)
	}

	return result, nil
}

// Up implements Service.Up. It builds the image if needed, creates a container
// and start it.
func (s *Service) Up() error {
	imageName, err := s.build()
	if err != nil {
		return err
	}

	return s.up(imageName, true)
}

// Info implements Service.Info. It returns an project.InfoSet with the containers
// related to this service (can be multiple if using the scale command).
func (s *Service) Info() (project.InfoSet, error) {
	result := project.InfoSet{}
	containers, err := s.collectContainers()
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		info, err := c.Info()
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

// Start implements Service.Start. It tries to start a container without creating it.
func (s *Service) Start() error {
	return s.up("", false)
}

func (s *Service) up(imageName string, create bool) error {
	containers, err := s.collectContainers()
	if err != nil {
		return err
	}

	logrus.Debugf("Found %d existing containers for service %s", len(containers), s.name)

	if len(containers) == 0 && create {
		c, err := s.createOne()
		if err != nil {
			return err
		}
		containers = []*Container{c}
	}

	return s.eachContainer(func(c *Container) error {
		if outOfSync, err := c.OutOfSync(); err != nil {
			return err
		} else if outOfSync {
			logrus.Warnf("%s needs rebuilding", s.Name())
		}
		return c.Up(imageName)
	})
}

func (s *Service) eachContainer(action func(*Container) error) error {
	containers, err := s.collectContainers()
	if err != nil {
		return err
	}

	tasks := utils.InParallel{}
	for _, container := range containers {
		task := func(container *Container) func() error {
			return func() error {
				return action(container)
			}
		}(container)

		tasks.Add(task)
	}

	return tasks.Wait()
}

// Down implements Service.Down. It stops any containers related to the service.
func (s *Service) Down() error {
	return s.eachContainer(func(c *Container) error {
		return c.Down()
	})
}

// Restart implements Service.Restart. It restarts any containers related to the service.
func (s *Service) Restart() error {
	return s.eachContainer(func(c *Container) error {
		return c.Restart()
	})
}

// Kill implements Service.Kill. It kills any containers related to the service.
func (s *Service) Kill() error {
	return s.eachContainer(func(c *Container) error {
		return c.Kill()
	})
}

// Delete implements Service.Delete. It removes any containers related to the service.
func (s *Service) Delete() error {
	return s.eachContainer(func(c *Container) error {
		return c.Delete()
	})
}

// Log implements Service.Log. It returns the docker logs for each container related to the service.
func (s *Service) Log() error {
	return s.eachContainer(func(c *Container) error {
		return c.Log()
	})
}

// Scale implements Service.Scale. It creates or removes containers to have the specified number
// of related container to the service to run.
func (s *Service) Scale(scale int) error {
	foundCount := 0
	err := s.eachContainer(func(c *Container) error {
		foundCount++
		if foundCount > scale {
			err := c.Down()
			if err != nil {
				return err
			}

			return c.Delete()
		}
		return nil
	})

	if err != nil {
		return err
	}

	if foundCount != scale {
		_, err := s.constructContainers(true, scale)
		if err != nil {
			return err
		}

	}

	return s.up("", false)
}

// Pull implements Service.Pull. It pulls or build the image of the service.
func (s *Service) Pull() error {
	containers, err := s.constructContainers(false, 1)
	if err != nil {
		return err
	}

	return containers[0].Pull()
}

// Containers implements Service.Containers. It returns the list of containers
// that are related to the service.
func (s *Service) Containers() ([]project.Container, error) {
	result := []project.Container{}
	containers, err := s.collectContainers()
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		result = append(result, c)
	}

	return result, nil
}
