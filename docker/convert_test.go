package docker

import (
	"path/filepath"
	"testing"

	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
	shlex "github.com/flynn/go-shlex"
	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	exp := []string{"sh", "-c", "exec /opt/bin/flanneld -logtostderr=true -iface=${NODE_IP}"}
	cmd, err := shlex.Split("sh -c 'exec /opt/bin/flanneld -logtostderr=true -iface=${NODE_IP}'")
	assert.Nil(t, err)
	assert.Equal(t, exp, cmd)
}

func TestParseBindsAndVolumes(t *testing.T) {
	ctx := &Context{}
	ctx.ComposeFile = "foo/docker-compose.yml"
	ctx.ResourceLookup = &lookup.FileConfigLookup{}

	abs, err := filepath.Abs(".")
	assert.Nil(t, err)
	cfg, hostCfg, err := Convert(&project.ServiceConfig{
		Volumes: []string{"/foo", "/home:/home", "/bar/baz", ".:/home", "/usr/lib:/usr/lib:ro"},
	}, ctx)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"/foo": {}, "/bar/baz": {}}, cfg.Volumes)
	assert.Equal(t, []string{"/home:/home", abs + "/foo:/home", "/usr/lib:/usr/lib:ro"}, hostCfg.Binds)
}
