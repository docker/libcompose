package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestScale(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	name2 := fmt.Sprintf("%s_%s_2", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)

	containers := s.GetContainersByProject(c, p)
	c.Assert(1, Equals, len(containers))

	s.FromText(c, p, "scale", "hello=2", SimpleTemplate)

	containers = s.GetContainersByProject(c, p)
	c.Assert(2, Equals, len(containers))

	for _, name := range []string{name, name2} {
		cn := s.GetContainerByName(c, name)
		c.Assert(cn, NotNil)
		c.Assert(cn.State.Running, Equals, true)
	}

	s.FromText(c, p, "scale", "--timeout", "0", "hello=1", SimpleTemplate)
	containers = s.GetContainersByProject(c, p)
	c.Assert(1, Equals, len(containers))

	cn = s.GetContainerByName(c, name2)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, IsNil)
	cn = s.GetContainerByName(c, name)
	c.Assert(cn, IsNil)
}
