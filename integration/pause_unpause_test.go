package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestPause(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, false)

	s.FromText(c, p, "pause", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, true)
}

func (s *RunSuite) TestPauseAlreadyPausedService(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, false)

	s.FromText(c, p, "pause", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, true)

	s.FromText(c, p, "pause", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, true)
}

func (s *RunSuite) TestUnpause(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, false)

	s.FromText(c, p, "pause", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, true)

	s.FromText(c, p, "unpause", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, false)
}

func (s *RunSuite) TestUnpauseNotPausedService(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, false)

	s.FromText(c, p, "unpause", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
	c.Assert(cn.State.Paused, Equals, false)
}
