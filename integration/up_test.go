package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestUp(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)
}

func (s *RunSuite) TestUpNotExistService(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "not_exist")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, IsNil)
}

func (s *RunSuite) TestLink(c *C) {
	p := s.ProjectFromText(c, "up", `
        server:
          image: busybox
          command: cat
          stdin_open: true
          expose:
          - 80
        client:
          image: busybox
          links:
          - server:foo
          - server
        `)

	serverName := fmt.Sprintf("%s_%s_1", p, "server")

	cn := s.GetContainerByName(c, serverName)
	c.Assert(cn, NotNil)
	c.Assert(cn.Config.ExposedPorts, DeepEquals, map[string]struct{}{
		"80/tcp": {},
	})

	clientName := fmt.Sprintf("%s_%s_1", p, "client")
	cn = s.GetContainerByName(c, clientName)
	c.Assert(cn, NotNil)
	c.Assert(asMap(cn.HostConfig.Links), DeepEquals, asMap([]string{
		fmt.Sprintf("/%s:/%s/%s", serverName, clientName, "foo"),
		fmt.Sprintf("/%s:/%s/%s", serverName, clientName, "server"),
		fmt.Sprintf("/%s:/%s/%s", serverName, clientName, serverName),
	}))
}
