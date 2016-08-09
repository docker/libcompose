package docker

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/labels"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/events"
	util "github.com/docker/libcompose/utils"
)

func (s *Service) createContainer(ctx context.Context, namer Namer, oldContainer string, configOverride *config.ServiceConfig, oneOff bool) (*Container, error) {
	serviceConfig := s.serviceConfig
	if configOverride != nil {
		serviceConfig.Command = configOverride.Command
		serviceConfig.Tty = configOverride.Tty
		serviceConfig.StdinOpen = configOverride.StdinOpen
	}
	configWrapper, err := ConvertToAPI(serviceConfig, s.context.Context, s.clientFactory)
	if err != nil {
		return nil, err
	}
	configWrapper.Config.Image = s.imageName()

	containerName, containerNumber := namer.Next()

	configWrapper.Config.Labels[labels.SERVICE.Str()] = s.name
	configWrapper.Config.Labels[labels.PROJECT.Str()] = s.project.Name
	configWrapper.Config.Labels[labels.HASH.Str()] = config.GetServiceHash(s.name, serviceConfig)
	configWrapper.Config.Labels[labels.ONEOFF.Str()] = strings.Title(strconv.FormatBool(oneOff))
	configWrapper.Config.Labels[labels.NUMBER.Str()] = fmt.Sprintf("%d", containerNumber)
	configWrapper.Config.Labels[labels.VERSION.Str()] = ComposeVersion

	err = s.populateAdditionalHostConfig(configWrapper.HostConfig)
	if err != nil {
		return nil, err
	}

	// FIXME(vdemeester): oldContainer should be a Container instead of a string
	client := s.clientFactory.Create(s)
	if oldContainer != "" {
		info, err := client.ContainerInspect(ctx, oldContainer)
		if err != nil {
			return nil, err
		}
		configWrapper.HostConfig.Binds = util.Merge(configWrapper.HostConfig.Binds, volumeBinds(configWrapper.Config.Volumes, &info))
	}

	logrus.Debugf("Creating container %s %#v", containerName, configWrapper)
	// FIXME(vdemeester): long-term will be container.Create(…)
	s.project.Notify(events.NewContainerCreateStartEvent(s.name, containerName))

	container, err := CreateContainer(ctx, client, containerName, configWrapper)
	if err != nil {
		s.project.Notify(events.NewContainerCreateFailedEvent(s.name, containerName, err))
		return nil, err
	}
	s.project.Notify(events.NewContainerCreateDoneEvent(s.name, containerName))

	return container, nil
}

func (s *Service) populateAdditionalHostConfig(hostConfig *container.HostConfig) error {
	links, err := s.getLinks()
	if err != nil {
		return err
	}

	for _, link := range s.DependentServices() {
		if !s.project.ServiceConfigs.Has(link.Target) {
			continue
		}

		service, err := s.project.CreateService(link.Target)
		if err != nil {
			return err
		}

		containers, err := service.Containers(context.Background())
		if err != nil {
			return err
		}

		if link.Type == project.RelTypeIpcNamespace {
			hostConfig, err = addIpc(hostConfig, service, containers, s.serviceConfig.Ipc)
		} else if link.Type == project.RelTypeNetNamespace {
			hostConfig, err = addNetNs(hostConfig, service, containers, s.serviceConfig.NetworkMode)
		}

		if err != nil {
			return err
		}
	}

	hostConfig.Links = []string{}
	for k, v := range links {
		hostConfig.Links = append(hostConfig.Links, strings.Join([]string{v, k}, ":"))
	}
	for _, v := range s.serviceConfig.ExternalLinks {
		hostConfig.Links = append(hostConfig.Links, v)
	}

	return nil
}

// FIXME(vdemeester) this is temporary
func (s *Service) getLinks() (map[string]string, error) {
	links := map[string]string{}
	for _, link := range s.DependentServices() {
		if !s.project.ServiceConfigs.Has(link.Target) {
			continue
		}

		service, err := s.project.CreateService(link.Target)
		if err != nil {
			return nil, err
		}

		// FIXME(vdemeester) container should not know service
		containers, err := service.Containers(context.Background())
		if err != nil {
			return nil, err
		}

		if link.Type == project.RelTypeLink {
			addLinks(links, service, link, containers)
		}

		if err != nil {
			return nil, err
		}
	}
	return links, nil
}

func addLinks(links map[string]string, service project.Service, rel project.ServiceRelationship, containers []project.Container) {
	for _, container := range containers {
		if _, ok := links[rel.Alias]; !ok {
			links[rel.Alias] = container.Name()
		}

		links[container.Name()] = container.Name()
	}
}

func addIpc(config *container.HostConfig, service project.Service, containers []project.Container, ipc string) (*container.HostConfig, error) {
	if len(containers) == 0 {
		return nil, fmt.Errorf("Failed to find container for IPC %v", ipc)
	}

	id, err := containers[0].ID()
	if err != nil {
		return nil, err
	}

	config.IpcMode = container.IpcMode("container:" + id)
	return config, nil
}

func addNetNs(config *container.HostConfig, service project.Service, containers []project.Container, networkMode string) (*container.HostConfig, error) {
	if len(containers) == 0 {
		return nil, fmt.Errorf("Failed to find container for networks ns %v", networkMode)
	}

	id, err := containers[0].ID()
	if err != nil {
		return nil, err
	}

	config.NetworkMode = container.NetworkMode("container:" + id)
	return config, nil
}

func volumeBinds(volumes map[string]struct{}, container *types.ContainerJSON) []string {
	result := make([]string, 0, len(container.Mounts))
	for _, mount := range container.Mounts {
		if _, ok := volumes[mount.Destination]; ok {
			result = append(result, fmt.Sprint(mount.Source, ":", mount.Destination))
		}
	}
	return result
}
