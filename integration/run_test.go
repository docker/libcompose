package integration

import (
	"fmt"
	"os/exec"

	. "gopkg.in/check.v1"
)

// FIXME find out why it fails with "inappropriate ioctl for device"
func (s *RunSuite) TestRun(c *C) {
	p := s.RandomProject()
	cmd := exec.Command(s.command, "-f", "./assets/run/docker-compose.yml", "-p", p, "run", "hello", "ls")

	output, err := cmd.CombinedOutput()
	c.Assert(err, IsNil, Commentf("%s", output))

	name := fmt.Sprintf("%s_%s_run_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, false)
	c.Assert(string(output), Equals, "test")
}
