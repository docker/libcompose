package auth

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/config/configfile"
	"github.com/docker/docker/registry"
)

// Lookup defines a method for looking up authentication information
type Lookup interface {
	All() map[string]types.AuthConfig
	Lookup(repoInfo *registry.RepositoryInfo) types.AuthConfig
	GetAuthConfigMap() map[string]Config
	SetAuthConfigMap(configMap map[string]Config) error
}

// ConfigLookup implements AuthLookup by reading a Docker config file
type ConfigLookup struct {
	*configfile.ConfigFile
}

// Config contains authorization information for connecting to a Registry
type Config struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	// Email is an optional value associated with the username.
	// This field is deprecated and will be removed in a later
	// version of docker.
	Email string `json:"email,omitempty"`

	ServerAddress string `json:"serveraddress,omitempty"`

	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `json:"identitytoken,omitempty"`

	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `json:"registrytoken,omitempty"`
}

// SetAuthConfigMap Update the docker auth config maps
func (c *ConfigLookup) SetAuthConfigMap(configMap map[string]Config) error {
	if c.ConfigFile == nil {
		return errors.New("ConfigFile not set in lookup")
	}
	authConfigMap := map[string]types.AuthConfig{}
	for k, v := range configMap {
		authConfigMap[k] = configToDockerTypeAuth(v)
	}
	c.ConfigFile.AuthConfigs = authConfigMap
	return nil
}

// GetAuthConfigMap This will return all of the AuthConfigurations
func (c *ConfigLookup) GetAuthConfigMap() map[string]Config {
	if c.ConfigFile == nil {
		return map[string]Config{}
	}
	configMap := map[string]Config{}
	for k, v := range c.All() {
		configMap[k] = dockerTypeAuthToConfig(v)
	}
	return configMap
}

func configToDockerTypeAuth(c Config) types.AuthConfig {
	return types.AuthConfig{
		Username:      c.Username,
		Password:      c.Password,
		Auth:          c.Auth,
		Email:         c.Email,
		ServerAddress: c.ServerAddress,
		IdentityToken: c.IdentityToken,
		RegistryToken: c.RegistryToken,
	}
}

func dockerTypeAuthToConfig(ac types.AuthConfig) Config {
	return Config{
		Username:      ac.Username,
		Password:      ac.Password,
		Auth:          ac.Auth,
		Email:         ac.Email,
		ServerAddress: ac.ServerAddress,
		IdentityToken: ac.IdentityToken,
		RegistryToken: ac.RegistryToken,
	}
}

// NewConfigLookup creates a new ConfigLookup for a given context
func NewConfigLookup(configfile *configfile.ConfigFile) *ConfigLookup {
	return &ConfigLookup{
		ConfigFile: configfile,
	}
}

// Lookup uses a Docker config file to lookup authentication information
func (c *ConfigLookup) Lookup(repoInfo *registry.RepositoryInfo) types.AuthConfig {
	if c.ConfigFile == nil || repoInfo == nil || repoInfo.Index == nil {
		return types.AuthConfig{}
	}
	return registry.ResolveAuthConfig(c.ConfigFile.AuthConfigs, repoInfo.Index)
}

// All uses a Docker config file to get all authentication information
func (c *ConfigLookup) All() map[string]types.AuthConfig {
	if c.ConfigFile == nil {
		return map[string]types.AuthConfig{}
	}
	return c.ConfigFile.AuthConfigs
}
