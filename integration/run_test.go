package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestRun(c *C) {
	projectName := s.RandomProject()
	p := s.FromText(c, projectName, "run", "hello", "echo", "test", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_run_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, false)
}
