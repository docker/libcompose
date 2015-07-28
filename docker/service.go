package docker

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/utils"
)

type Service struct {
	name          string
	serviceConfig *project.ServiceConfig
	context       *Context
	imageName     string
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) Config() *project.ServiceConfig {
	return s.serviceConfig
}

func (s *Service) DependentServices() []project.ServiceRelationship {
	return project.DefaultDependentServices(s.context.Project, s)
}

func (s *Service) Create() error {
	_, err := s.createOne()
	return err
}

func (s *Service) collectContainers() ([]*Container, error) {
	containers, err := GetContainersByFilter(s.context.Client, SERVICE.Eq(s.name), PROJECT.Eq(s.context.Project.Name))
	if err != nil {
		return nil, err
	}

	result := []*Container{}

	for _, container := range containers {
		result = append(result, NewContainer(s.context.Client, container.Labels[NAME.Str()], s))
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

	namer := NewNamer(s.context.Client, s.context.Project.Name, s.name)
	defer namer.Close()

	for i := len(result); i < count; i++ {
		containerName := namer.Next()

		c := NewContainer(s.context.Client, containerName, s)

		if create {
			imageName, err := s.build()
			if err != nil {
				return nil, err
			}

			dockerContainer, err := c.Create(imageName)
			if err != nil {
				return nil, err
			} else {
				logrus.Debugf("Created container %s: %v", dockerContainer.Id, dockerContainer.Names)
			}
		}

		result = append(result, NewContainer(s.context.Client, containerName, s))
	}

	return result, nil
}

func (s *Service) Up() error {
	imageName, err := s.build()
	if err != nil {
		return err
	}

	return s.up(imageName, true)
}

func (s *Service) Info() (project.InfoSet, error) {
	result := project.InfoSet{}
	containers, err := s.collectContainers()
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if info, err := c.Info(); err != nil {
			return nil, err
		} else {
			result = append(result, info)
		}
	}

	return result, nil
}

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

func (s *Service) Down() error {
	return s.eachContainer(func(c *Container) error {
		return c.Down()
	})
}

func (s *Service) Restart() error {
	return s.eachContainer(func(c *Container) error {
		return c.Restart()
	})
}

func (s *Service) Kill() error {
	return s.eachContainer(func(c *Container) error {
		return c.Kill()
	})
}

func (s *Service) Delete() error {
	return s.eachContainer(func(c *Container) error {
		return c.Delete()
	})
}

func (s *Service) Log() error {
	return s.eachContainer(func(c *Container) error {
		return c.Log()
	})
}

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

func (s *Service) Pull() error {
	containers, err := s.constructContainers(false, 1)
	if err != nil {
		return err
	}

	return containers[0].Pull()
}

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
