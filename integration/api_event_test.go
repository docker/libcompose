package integration

import (
	"time"

	"golang.org/x/net/context"
	. "gopkg.in/check.v1"

	eventtypes "github.com/docker/engine-api/types/events"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

func (s *APISuite) TestEvents(c *C) {
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
	c.Assert(err, IsNil)

	ctx, cancelFun := context.WithCancel(context.Background())

	messages, err := project.Events(ctx)
	c.Assert(err, IsNil)

	go func() {
		c.Assert(project.Up(ctx, options.Up{}), IsNil)
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
	c.Assert(len(events), Equals, 4, Commentf("%v", events))
}
