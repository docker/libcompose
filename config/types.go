package config

import (
	"sync"

	"github.com/docker/libcompose/yaml"
)

// EnvironmentLookup defines methods to provides environment variable loading.
type EnvironmentLookup interface {
	Lookup(key, serviceName string, config *ServiceConfig) []string
}

// ResourceLookup defines methods to provides file loading.
type ResourceLookup interface {
	Lookup(file, relativeTo string) ([]byte, string, error)
	ResolvePath(path, inFile string) string
}

// ServiceConfigV1 holds version 1 of libcompose service configuration
type ServiceConfigV1 struct {
	Build         string               `yaml:"build,omitempty"`
	CapAdd        []string             `yaml:"cap_add,omitempty"`
	CapDrop       []string             `yaml:"cap_drop,omitempty"`
	CgroupParent  string               `yaml:"cgroup_parent,omitempty"`
	CPUQuota      int64                `yaml:"cpu_quota,omitempty"`
	CPUSet        string               `yaml:"cpuset,omitempty"`
	CPUShares     int64                `yaml:"cpu_shares,omitempty"`
	Command       yaml.Command         `yaml:"command,flow,omitempty"`
	ContainerName string               `yaml:"container_name,omitempty"`
	Devices       []string             `yaml:"devices,omitempty"`
	DNS           yaml.Stringorslice   `yaml:"dns,omitempty"`
	DNSSearch     yaml.Stringorslice   `yaml:"dns_search,omitempty"`
	Dockerfile    string               `yaml:"dockerfile,omitempty"`
	DomainName    string               `yaml:"domainname,omitempty"`
	Entrypoint    yaml.Command         `yaml:"entrypoint,flow,omitempty"`
	EnvFile       yaml.Stringorslice   `yaml:"env_file,omitempty"`
	Environment   yaml.MaporEqualSlice `yaml:"environment,omitempty"`
	Hostname      string               `yaml:"hostname,omitempty"`
	Image         string               `yaml:"image,omitempty"`
	Labels        yaml.SliceorMap      `yaml:"labels,omitempty"`
	Links         yaml.MaporColonSlice `yaml:"links,omitempty"`
	LogDriver     string               `yaml:"log_driver,omitempty"`
	MacAddress    string               `yaml:"mac_address,omitempty"`
	MemLimit      int64                `yaml:"mem_limit,omitempty"`
	MemSwapLimit  int64                `yaml:"memswap_limit,omitempty"`
	Name          string               `yaml:"name,omitempty"`
	Net           string               `yaml:"net,omitempty"`
	Pid           string               `yaml:"pid,omitempty"`
	Uts           string               `yaml:"uts,omitempty"`
	Ipc           string               `yaml:"ipc,omitempty"`
	Ports         []string             `yaml:"ports,omitempty"`
	Privileged    bool                 `yaml:"privileged,omitempty"`
	Restart       string               `yaml:"restart,omitempty"`
	ReadOnly      bool                 `yaml:"read_only,omitempty"`
	StdinOpen     bool                 `yaml:"stdin_open,omitempty"`
	SecurityOpt   []string             `yaml:"security_opt,omitempty"`
	Tty           bool                 `yaml:"tty,omitempty"`
	User          string               `yaml:"user,omitempty"`
	VolumeDriver  string               `yaml:"volume_driver,omitempty"`
	Volumes       []string             `yaml:"volumes,omitempty"`
	VolumesFrom   []string             `yaml:"volumes_from,omitempty"`
	WorkingDir    string               `yaml:"working_dir,omitempty"`
	Expose        []string             `yaml:"expose,omitempty"`
	ExternalLinks []string             `yaml:"external_links,omitempty"`
	LogOpt        map[string]string    `yaml:"log_opt,omitempty"`
	ExtraHosts    []string             `yaml:"extra_hosts,omitempty"`
	Ulimits       yaml.Ulimits         `yaml:"ulimits,omitempty"`
}

// Build holds v2 build information
type Build struct {
	Context    string               `yaml:"context,omitempty"`
	Dockerfile string               `yaml:"dockerfile,omitempty"`
	Args       yaml.MaporEqualSlice `yaml:"args,omitempty"`
}

// Log holds v2 logging information
type Log struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

// ServiceConfig holds version 2 of libcompose service configuration
type ServiceConfig struct {
	Build         Build                `yaml:"build,omitempty"`
	CapAdd        []string             `yaml:"cap_add,omitempty"`
	CapDrop       []string             `yaml:"cap_drop,omitempty"`
	CPUSet        string               `yaml:"cpuset,omitempty"`
	CPUShares     int64                `yaml:"cpu_shares,omitempty"`
	CPUQuota      int64                `yaml:"cpu_quota,omitempty"`
	Command       yaml.Command         `yaml:"command,flow,omitempty"`
	CgroupParent  string               `yaml:"cgroup_parrent,omitempty"`
	ContainerName string               `yaml:"container_name,omitempty"`
	Devices       []string             `yaml:"devices,omitempty"`
	DependsOn     []string             `yaml:"depends_on,omitempty"`
	DNS           yaml.Stringorslice   `yaml:"dns,omitempty"`
	DNSSearch     yaml.Stringorslice   `yaml:"dns_search,omitempty"`
	DomainName    string               `yaml:"domain_name,omitempty"`
	Entrypoint    yaml.Command         `yaml:"entrypoint,flow,omitempty"`
	EnvFile       yaml.Stringorslice   `yaml:"env_file,omitempty"`
	Environment   yaml.MaporEqualSlice `yaml:"environment,omitempty"`
	Expose        []string             `yaml:"expose,omitempty"`
	Extends       yaml.MaporEqualSlice `yaml:"extends,omitempty"`
	ExternalLinks []string             `yaml:"external_links,omitempty"`
	ExtraHosts    []string             `yaml:"extra_hosts,omitempty"`
	Image         string               `yaml:"image,omitempty"`
	Hostname      string               `yaml:"hostname,omitempty"`
	Ipc           string               `yaml:"ipc,omitempty"`
	Labels        yaml.SliceorMap      `yaml:"labels,omitempty"`
	Links         yaml.MaporColonSlice `yaml:"links,omitempty"`
	Logging       Log                  `yaml:"logging,omitempty"`
	MacAddress    string               `yaml:"mac_address,omitempty"`
	MemLimit      int64                `yaml:"mem_limit,omitempty"`
	MemSwapLimit  int64                `yaml:"memswap_limit,omitempty"`
	NetworkMode   string               `yaml:"network_mode,omitempty"`
	Networks      []string             `yaml:"networks,omitempty"`
	Pid           string               `yaml:"pid,omitempty"`
	Ports         []string             `yaml:"ports,omitempty"`
	Privileged    bool                 `yaml:"privileged,omitempty"`
	SecurityOpt   []string             `yaml:"security_opt,omitempty"`
	StopSignal    string               `yaml:"stop_signal,omitempty"`
	VolumeDriver  string               `yaml:"volume_driver,omitempty"`
	Volumes       []string             `yaml:"volumes,omitempty"`
	VolumesFrom   []string             `yaml:"volumes_from,omitempty"`
	Uts           string               `yaml:"uts,omitempty"`
	Restart       string               `yaml:"restart,omitempty"`
	ReadOnly      bool                 `yaml:"read_only,omitempty"`
	StdinOpen     bool                 `yaml:"stdin_open,omitempty"`
	Tty           bool                 `yaml:"tty,omitempty"`
	User          string               `yaml:"user,omitempty"`
	WorkingDir    string               `yaml:"working_dir,omitempty"`
	Ulimits       yaml.Ulimits         `yaml:"ulimits,omitempty"`
}

// VolumeConfig holds v2 volume configuration
type VolumeConfig struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	External   bool              `yaml:"external,omitempty"`
}

// Ipam holds v2 network IPAM information
type Ipam struct {
	Driver string   `yaml:"driver,omitempty"`
	Config []string `yaml:"config,omitempty"`
}

// NetworkConfig holds v2 network configuration
type NetworkConfig struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	External   bool              `yaml:"external,omitempty"`
	Ipam       Ipam              `yaml:"ipam,omitempty"`
}

// Config holds libcompose top level configuration
type Config struct {
	Version  string                    `yaml:"version,omitempty"`
	Services RawServiceMap             `yaml:"services,omitempty"`
	Volumes  map[string]*VolumeConfig  `yaml:"volumes,omitempty"`
	Networks map[string]*NetworkConfig `yaml:"networks,omitempty"`
}

// NewServiceConfigs initializes a new Configs struct
func NewServiceConfigs() *ServiceConfigs {
	return &ServiceConfigs{
		m: make(map[string]*ServiceConfig),
	}
}

// ServiceConfigs holds a concurrent safe map of ServiceConfig
type ServiceConfigs struct {
	m  map[string]*ServiceConfig
	mu sync.RWMutex
}

// Has checks if the config map has the specified name
func (c *ServiceConfigs) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.m[name]
	return ok
}

// Get returns the config and the presence of the specified name
func (c *ServiceConfigs) Get(name string) (*ServiceConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	service, ok := c.m[name]
	return service, ok
}

// Add add the specifed config with the specified name
func (c *ServiceConfigs) Add(name string, service *ServiceConfig) {
	c.mu.Lock()
	c.m[name] = service
	c.mu.Unlock()
}

// Len returns the len of the configs
func (c *ServiceConfigs) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)
}

// Keys returns the names of the config
func (c *ServiceConfigs) Keys() []string {
	keys := []string{}
	c.mu.RLock()
	defer c.mu.RUnlock()
	for name := range c.m {
		keys = append(keys, name)
	}
	return keys
}

// RawService is represent a Service in map form unparsed
type RawService map[string]interface{}

// RawServiceMap is a collection of RawServices
type RawServiceMap map[string]RawService

// ParseOptions are a set of options to customize the parsing process
type ParseOptions struct {
	Interpolate bool
	Validate    bool
	Preprocess  func(RawServiceMap) (RawServiceMap, error)
	Postprocess func(map[string]*ServiceConfig) (map[string]*ServiceConfig, error)
}
