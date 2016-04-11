package integration

import (
	. "gopkg.in/check.v1"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

func (s *RunSuite) TestVolumeWithoutComposeFile(c *C) {
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

	err = project.Up()
	c.Assert(err, IsNil)
}
