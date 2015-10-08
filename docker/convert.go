package docker

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/nat"
	"github.com/docker/docker/runconfig"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/utils"
	"github.com/samalba/dockerclient"
)

// Filter filters the specified string slice with the specified function.
func Filter(vs []string, f func(string) bool) []string {
	r := make([]string, 0, len(vs))
	for _, v := range vs {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func isBind(s string) bool {
	return strings.ContainsRune(s, ':')
}

func isVolume(s string) bool {
	return !isBind(s)
}

// ConvertToAPI converts a service configuration to a docker API container configuration.
func ConvertToAPI(s *Service) (*dockerclient.ContainerConfig, error) {
	config, hostConfig, err := Convert(s.serviceConfig, s.context)
	if err != nil {
		return nil, err
	}

	var result dockerclient.ContainerConfig
	err = utils.ConvertByJSON(config, &result)
	if err != nil {
		logrus.Errorf("Failed to convert config to API structure: %v\n%#v", err, config)
		return nil, err
	}

	err = utils.ConvertByJSON(hostConfig, &result.HostConfig)
	if err != nil {
		logrus.Errorf("Failed to convert hostConfig to API structure: %v\n%#v", err, hostConfig)
	}
	return &result, err
}

// Convert converts a service configuration to an docker inner representation (using runconfig structures)
func Convert(c *project.ServiceConfig, ctx *Context) (*runconfig.Config, *runconfig.HostConfig, error) {
	volumes := make(map[string]struct{}, len(c.Volumes))
	for k, v := range c.Volumes {
		vol := ctx.ResourceLookup.ResolvePath(v, ctx.ComposeFile)

		c.Volumes[k] = vol
		if isVolume(vol) {
			volumes[vol] = struct{}{}
		}
	}

	ports, binding, err := nat.ParsePortSpecs(c.Ports)
	if err != nil {
		return nil, nil, err
	}
	restart, err := runconfig.ParseRestartPolicy(c.Restart)
	if err != nil {
		return nil, nil, err
	}

	exposedPorts, _, err := nat.ParsePortSpecs(c.Expose)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range exposedPorts {
		ports[k] = v
	}

	deviceMappings, err := parseDevices(c.Devices)
	if err != nil {
		return nil, nil, err
	}

	config := &runconfig.Config{
		Entrypoint:   runconfig.NewEntrypoint(c.Entrypoint.Slice()...),
		Hostname:     c.Hostname,
		Domainname:   c.DomainName,
		User:         c.User,
		Env:          c.Environment.Slice(),
		Cmd:          runconfig.NewCommand(c.Command.Slice()...),
		Image:        c.Image,
		Labels:       c.Labels.MapParts(),
		ExposedPorts: ports,
		Tty:          c.Tty,
		OpenStdin:    c.StdinOpen,
		WorkingDir:   c.WorkingDir,
		VolumeDriver: c.VolumeDriver,
		Volumes:      volumes,
	}
	hostConfig := &runconfig.HostConfig{
		VolumesFrom: c.VolumesFrom,
		CapAdd:      runconfig.NewCapList(c.CapAdd),
		CapDrop:     runconfig.NewCapList(c.CapDrop),
		CPUShares:   c.CPUShares,
		CpusetCpus:  c.CPUSet,
		ExtraHosts:  c.ExtraHosts,
		Privileged:  c.Privileged,
		Binds:       Filter(c.Volumes, isBind),
		Devices:     deviceMappings,
		DNS:         c.DNS.Slice(),
		DNSSearch:   c.DNSSearch.Slice(),
		LogConfig: runconfig.LogConfig{
			Type:   c.LogDriver,
			Config: c.LogOpt,
		},
		Memory:         c.MemLimit,
		MemorySwap:     c.MemSwapLimit,
		NetworkMode:    runconfig.NetworkMode(c.Net),
		ReadonlyRootfs: c.ReadOnly,
		PidMode:        runconfig.PidMode(c.Pid),
		UTSMode:        runconfig.UTSMode(c.Uts),
		IpcMode:        runconfig.IpcMode(c.Ipc),
		PortBindings:   binding,
		RestartPolicy:  restart,
		SecurityOpt:    c.SecurityOpt,
	}

	return config, hostConfig, nil
}

func parseDevices(devices []string) ([]runconfig.DeviceMapping, error) {
	// parse device mappings
	deviceMappings := []runconfig.DeviceMapping{}
	for _, device := range devices {
		deviceMapping, err := runconfig.ParseDevice(device)
		if err != nil {
			return nil, err
		}
		deviceMappings = append(deviceMappings, deviceMapping)
	}

	return deviceMappings, nil
}
