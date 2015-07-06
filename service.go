package libcompose

import ()

type Service struct {
	Name                 string
	LogPrefix            string
	EnvironmentVariables []string
	Env                  interface{} `yaml:"environment"`
	Expose               string      `yaml:"expose"`
	Image                string      `yaml:"image"`
	BuildDir             string      `yaml:"build"`
	Dns                  []string    `yaml:"dns"`
	NetworkingMode       string      `yaml:"net"`
	Command              string      `yaml:"command"`
	Links                []string    `yaml:"links"`
	Ports                []string    `yaml:"ports"`
	Volumes              []string    `yaml:"volumes"` // Volumes map[string]struct{} `yaml:"volumes"` // previous type for use with yaml version 1 : "gopkg.in/yaml.v1"
	VolumesFrom          []string    `yaml:"volumes_from"`
	WorkingDir           string      `yaml:"working_dir"`
	Entrypoint           string      `yaml:"entrypoint"`
	User                 string      `yaml:"user"`
	HostName             string      `yaml:"hostname"`
	DomainName           string      `yaml:"domainname"`
	MemLimit             string      `yaml:"mem_limit"`
	Privileged           bool        `yaml:"privileged"`
	WatchDirs            []string    `yaml:"watch"`
	//ExposedPorts map[apiClient.Port]struct{}
	//Container    apiClient.Container
	//OnAttachHook func(io.Reader, Service)
	//Api          *apiClient.Client
	//Cli          *dockerCli.DockerCli
}

// returns whether the Service has at least one dependency
func (s *Service) HasDependencies() bool {
	return len(s.Links) > 0 || len(s.VolumesFrom) > 0
}

// returns whether the Service has no dependencies
func (s *Service) HasNoDependencies() bool {
	return s.HasDependencies() == false
}

// returns the services' names on which the service is dependent
func (s *Service) GetDependenciesNames() []string {

	var dependencies []string

	// links
	for _, name := range s.Links {
		dependencies = addUniqueStringToSlice(name, dependencies)
	}

	// volumes_from
	for _, name := range s.VolumesFrom {
		dependencies = addUniqueStringToSlice(name, dependencies)
	}

	return dependencies
}

// returns whether all dependencies are part of the "services" slice parameter
// this method can be use for dependency fullfillment check
func (s *Service) AreDependenciesIn(services []*Service) bool {

	dependenciesNames := s.GetDependenciesNames()

	for _, d := range dependenciesNames {
		found := false
		for _, s := range services {
			if d == s.Name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

/////////////////////////////////////////////////////////////////////
// PRIVATE UTILITY FUNCTIONS
/////////////////////////////////////////////////////////////////////

// add a string to a []string slice, avoiding duplicates
// returns the slice containing uniqueStr
func addUniqueStringToSlice(uniqueStr string, slice []string) []string {
	for _, str := range slice {
		if str == uniqueStr {
			// uniqueStr is already in the slice
			return slice
		}
	}
	return append(slice, uniqueStr)
}
