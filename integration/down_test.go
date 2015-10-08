package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestDown(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")

	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)

	s.FromText(c, p, "down", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, false)
}
