package integration

import (
	"time"

	"golang.org/x/net/context"
	check "gopkg.in/check.v1"

	eventtypes "github.com/docker/engine-api/types/events"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

func (s *APISuite) TestEvents(c *check.C) {
	testRequires(c, not(DaemonVersionIs("1.9")))
	composeFile := `
simple:
  image: busybox:latest
  command: top
another:
  image: busybox:latest
  command: top
`
	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{[]byte(composeFile)},
			ProjectName:  "test-api-events",
		},
	}, nil)
	c.Assert(err, check.IsNil)

	ctx, cancelFun := context.WithCancel(context.Background())

	messages, err := project.Events(ctx)
	c.Assert(err, check.IsNil)

	go func() {
		c.Assert(project.Up(ctx, options.Up{}), check.IsNil)
		// Close after everything is done
		time.Sleep(250 * time.Millisecond)
		cancelFun()
		close(messages)
	}()

	events := []eventtypes.Message{}
	for message := range messages {
		events = append(events, message)
	}

	// Should be 4 events (2 create, 2 start)
	c.Assert(len(events), check.Equals, 4, check.Commentf("%v", events))
}
