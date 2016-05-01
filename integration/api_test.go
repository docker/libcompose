package integration

import (
	. "gopkg.in/check.v1"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

func init() {
	Suite(&APISuite{})
}

type APISuite struct{}

func (s *APISuite) TestVolumeWithoutComposeFile(c *C) {
	service := `
service:
  image: busybox
  command: echo Hello world!
  volumes:
    - /etc/selinux:/etc/selinux`

	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{[]byte(service)},
			ProjectName:  "test-volume-without-compose-file",
		},
	})

	c.Assert(err, IsNil)

	err = project.Up(options.Up{})
	c.Assert(err, IsNil)
}
