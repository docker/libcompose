package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestDelete(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")

	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)

	s.FromText(c, p, "rm", "--force", `
        hello:
          image: busybox
          stdin_open: true
          tty: true
        `)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, IsNil)
}
