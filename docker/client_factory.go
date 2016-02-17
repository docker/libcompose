package docker

import (
	"github.com/docker/engine-api/client"
	"github.com/docker/libcompose/project"
)

// ClientFactory is a factory to create docker clients.
type ClientFactory interface {
	// Create constructs a Docker client for the given service. The passed in
	// config may be nil in which case a generic client for the project should
	// be returned.
	Create(service project.Service) client.APIClient
}

type defaultClientFactory struct {
	client client.APIClient
}

// NewDefaultClientFactory creates and returns the default client factory that uses
// github.com/docker/engine-api client.
func NewDefaultClientFactory(opts ClientOpts) (ClientFactory, error) {
	client, err := CreateClient(opts)
	if err != nil {
		return nil, err
	}

	return &defaultClientFactory{
		client: client,
	}, nil
}

func (s *defaultClientFactory) Create(service project.Service) client.APIClient {
	return s.client
}
