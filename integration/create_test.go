package integration

import (
	"fmt"
	"os"

	. "gopkg.in/check.v1"
)

func (s *RunSuite) TestFields(c *C) {
	p := s.CreateProjectFromText(c, `
        hello:
          image: tianon/true
          cpuset: 1,2
          mem_limit: 4194304
        `)

	name := fmt.Sprintf("%s_%s_1", p, "hello")
	cn := s.GetContainerByName(c, name)
	c.Assert(cn, NotNil)

	c.Assert(cn.Config.Image, Equals, "tianon/true")
	c.Assert(cn.HostConfig.CPUSetCPUs, Equals, "1,2")
	c.Assert(cn.HostConfig.Memory, Equals, int64(4194304))
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

func (s *RunSuite) TestContainerName(c *C) {
	containerName := "containerName"
	template := fmt.Sprintf(`hello:
    image: busybox
    command: top
    container_name: %s`, containerName)
	s.CreateProjectFromText(c, template)

	cn := s.GetContainerByName(c, containerName)
	c.Assert(cn, NotNil)

	c.Assert(cn.Name, Equals, "/"+containerName)
}

func (s *RunSuite) TestContainerNameWithScale(c *C) {
	containerName := "containerName"
	template := fmt.Sprintf(`hello:
    image: busybox
    command: top
    container_name: %s`, containerName)
	p := s.CreateProjectFromText(c, template)

	s.FromText(c, p, "scale", "hello=2", template)
	containers := s.GetContainersByProject(c, p)
	c.Assert(len(containers), Equals, 1)

}

func (s *RunSuite) TestInterpolation(c *C) {
	os.Setenv("IMAGE", "tianon/true")

	p := s.CreateProjectFromText(c, `
        test:
          image: $IMAGE
        `)

	name := fmt.Sprintf("%s_%s_1", p, "test")
	testContainer := s.GetContainerByName(c, name)

	p = s.CreateProjectFromText(c, `
        reference:
          image: tianon/true
        `)

	name = fmt.Sprintf("%s_%s_1", p, "reference")
	referenceContainer := s.GetContainerByName(c, name)

	c.Assert(testContainer, NotNil)

	c.Assert(referenceContainer.Image, Equals, testContainer.Image)

	os.Unsetenv("IMAGE")
}

func (s *RunSuite) TestInterpolationWithExtends(c *C) {
	os.Setenv("IMAGE", "tianon/true")
	os.Setenv("TEST_PORT", "8000")

	p := s.CreateProjectFromText(c, `
        test:
                extends:
                        file: ./assets/interpolation/docker-compose.yml
                        service: base
                ports:
                        - ${TEST_PORT}
        `)

	name := fmt.Sprintf("%s_%s_1", p, "test")
	testContainer := s.GetContainerByName(c, name)

	p = s.CreateProjectFromText(c, `
	reference:
	  image: tianon/true
		ports:
		  - 8000
	`)

	name = fmt.Sprintf("%s_%s_1", p, "reference")
	referenceContainer := s.GetContainerByName(c, name)

	c.Assert(testContainer, NotNil)

	c.Assert(referenceContainer.Image, Equals, testContainer.Image)

	os.Unsetenv("TEST_PORT")
	os.Unsetenv("IMAGE")
}

func (s *RunSuite) TestFieldTypeConversions(c *C) {
	os.Setenv("LIMIT", "40000000")

	p := s.CreateProjectFromText(c, `
        test:
          image: tianon/true
          mem_limit: $LIMIT
          memswap_limit: "40000000"
          hostname: 100
        `)

	name := fmt.Sprintf("%s_%s_1", p, "test")
	testContainer := s.GetContainerByName(c, name)

	p = s.CreateProjectFromText(c, `
        reference:
          image: tianon/true
          mem_limit: 40000000
          memswap_limit: 40000000
          hostname: "100"
        `)

	name = fmt.Sprintf("%s_%s_1", p, "reference")
	referenceContainer := s.GetContainerByName(c, name)

	c.Assert(testContainer, NotNil)

	c.Assert(referenceContainer.Image, Equals, testContainer.Image)

	os.Unsetenv("LIMIT")
}
