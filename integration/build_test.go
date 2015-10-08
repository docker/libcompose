package integration

import (
	"fmt"
	"os"
	"os/exec"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestBuild(c *C) {
	p := s.RandomProject()
	cmd := exec.Command(s.command, "-f", "./assets/build/docker-compose.yml", "-p", p, "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	oneImageName := fmt.Sprintf("%s_one", p)
	twoImageName := fmt.Sprintf("%s_two", p)

	c.Assert(err, IsNil)

	client := GetClient(c)
	one, err := client.InspectImage(oneImageName)
	c.Assert(err, IsNil)
	c.Assert(one.Config.Cmd, DeepEquals, []string{"echo", "one"})

	two, err := client.InspectImage(twoImageName)
	c.Assert(err, IsNil)
	c.Assert(two.Config.Cmd, DeepEquals, []string{"echo", "two"})
}
