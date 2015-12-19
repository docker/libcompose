package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func (s *RunSuite) TestBuildWithNoCache1(c *C) {
	p := s.RandomProject()
	cmd := exec.Command(s.command, "-f", "./assets/build/docker-compose.yml", "-p", p, "build")

	output, err := cmd.Output()
	c.Assert(err, IsNil)

	cmd = exec.Command(s.command, "-f", "./assets/build/docker-compose.yml", "-p", p, "build")
	output, err = cmd.Output()
	c.Assert(err, IsNil)
	out := string(output[:])
	c.Assert(strings.Contains(out,
		"Using cache"),
		Equals, true, Commentf("%s", out))
}

func (s *RunSuite) TestBuildWithNoCache2(c *C) {
	p := s.RandomProject()
	cmd := exec.Command(s.command, "-f", "./assets/build/docker-compose.yml", "-p", p, "build")

	output, err := cmd.Output()
	c.Assert(err, IsNil)

	cmd = exec.Command(s.command, "-f", "./assets/build/docker-compose.yml", "-p", p, "build", "--no-cache")
	output, err = cmd.Output()
	c.Assert(err, IsNil)
	out := string(output[:])
	c.Assert(strings.Contains(out,
		"Using cache"),
		Equals, false, Commentf("%s", out))
}

func (s *RunSuite) TestBuildWithNoCache3(c *C) {
	p := s.RandomProject()
	cmd := exec.Command(s.command, "-f", "./assets/build/docker-compose.yml", "-p", p, "build", "--no-cache")
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
