package docker

import (
	dockerclient "github.com/fsouza/go-dockerclient"
)

// GetContainersByFilter looks up the hosts containers with the specified filters and
// returns a list of container matching it, or an error.
func GetContainersByFilter(client *dockerclient.Client, filters ...map[string][]string) ([]dockerclient.APIContainers, error) {
	var filterResult map[string][]string

	for _, filter := range filters {
		if filterResult == nil {
			filterResult = filter
		} else {
			filterResult = And(filterResult, filter)
		}
	}

	return client.ListContainers(dockerclient.ListContainersOptions{All: true, Filters: filterResult})
}

// GetContainerByName looks up the hosts containers with the specified name and
// returns it, or an error.
func GetContainerByName(client *dockerclient.Client, name string) (*dockerclient.APIContainers, error) {
	containers, err := client.ListContainers(dockerclient.ListContainersOptions{All: true, Filters: NAME.Eq(name)})
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func GetContainerById(client dockerclient.Client, id string) (*dockerclient.Container, error) {
	containers, err := client.ListContainers(true, false, "")
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, nil
	}

	for _, c := range containers {
		if c.Id == id {
			return &c, nil
		}
	}

	return nil, nil
}
