package project

import (
	"github.com/docker/libcompose/project/events"
	"github.com/docker/libcompose/project/options"
)

// APIProject is an interface defining the methods a libcompose project should implement.
type APIProject interface {
	Build(options options.Build, sevice ...string) error
	Create(options options.Create, services ...string) error
	Delete(options options.Delete, services ...string) error
	Down(options options.Down, services ...string) error
	Kill(signal string, services ...string) error
	Log(follow bool, services ...string) error
	Pause(services ...string) error
	Ps(onlyID bool, services ...string) (InfoSet, error)
	// FIXME(vdemeester) we could use nat.Port instead ?
	Port(index int, protocol, serviceName, privatePort string) (string, error)
	Pull(services ...string) error
	Restart(timeout int, services ...string) error
	Run(serviceName string, commandParts []string) (int, error)
	Scale(timeout int, servicesScale map[string]int) error
	Start(services ...string) error
	Stop(timeout int, services ...string) error
	Unpause(services ...string) error
	Up(options options.Up, services ...string) error

	Parse() error

	// FIXME(vdemeester) move outside this interface. Delete(…) should take a
	// function for the force thingie (ask a question or not)
	ListStoppedContainers(services ...string) ([]string, error)

	// FIXME(vdemeester) should be moved outside this interface
	// The listener/notify mecanism could/should be indenpendant from Project
	AddListener(c chan<- events.Event)
	Notify(eventType events.EventType, serviceName string, data map[string]string)
}
