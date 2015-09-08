package integration

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "gopkg.in/check.v1"
)

const SimpleTemplate = `
	hello:
	  image: busybox
	  stdin_open: true
	  tty: true
	`

func Test(t *testing.T) { TestingT(t) }

func (s *RunSuite) TestFields(c *C) {
	p := s.CreateProjectFromText(c, `
	hello:
	  image: tianon/true
	  cpuset: 1,2
	  mem_limit: 4194304
	  memswap_limit: 8388608
	`)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.Config.Cpuset, Equals, "1,2")
	c.Assert(cn.HostConfig.Memory, Equals, int64(4194304))
	c.Assert(cn.HostConfig.MemorySwap, Equals, int64(8388608))
}

func (s *RunSuite) TestHelloWorld(c *C) {
	p := s.CreateProjectFromText(c, `
	hello:
	  image: tianon/true
	`)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.Name, Equals, "/"+name)
}

func (s *RunSuite) TestUp(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)
}

func (s *RunSuite) TestRestart(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.State.Running, Equals, true)
	time := cn.State.StartedAt.UnixNano()

	s.FromText(c, p, "restart", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)

	c.Assert(time, Not(Equals), cn.State.StartedAt.UnixNano())
}

func (s *RunSuite) TestStop(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")

	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)

	s.FromText(c, p, "stop", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, false)
}

func (s *RunSuite) TestKill(c *C) {
	p := s.ProjectFromText(c, "up", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")

	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)

	s.FromText(c, p, "kill", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, false)
}

func (s *RunSuite) TestStart(c *C) {
	p := s.ProjectFromText(c, "create", SimpleTemplate)

	name := fmt.Sprintf("%s_%s_1", p, "hello")

	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, false)

	s.FromText(c, p, "start", SimpleTemplate)

	cn = s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)
	c.Assert(cn.State.Running, Equals, true)
}

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
}

func (s *RunSuite) TestPull(c *C) {
	//TODO: This doesn't test much
	s.ProjectFromText(c, "pull", `
	hello:
	  image: tianon/true
	  stdin_open: true
	  tty: true
	`)
}

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

func asMap(items []string) map[string]bool {
	result := map[string]bool{}
	for _, item := range items {
		result[item] = true
	}
	return result
}
