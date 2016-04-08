package project

import (
	"fmt"
	"github.com/docker/libcompose/config"
)

// EventType defines a type of libcompose event.
type EventType int

// Definitions of libcompose events
const (
	NoEvent = EventType(iota)

	EventContainerCreated = EventType(iota)
	EventContainerStarted = EventType(iota)

	EventServiceAdd          = EventType(iota)
	EventServiceUpStart      = EventType(iota)
	EventServiceUpIgnored    = EventType(iota)
	EventServiceUp           = EventType(iota)
	EventServiceCreateStart  = EventType(iota)
	EventServiceCreate       = EventType(iota)
	EventServiceDeleteStart  = EventType(iota)
	EventServiceDelete       = EventType(iota)
	EventServiceDownStart    = EventType(iota)
	EventServiceDown         = EventType(iota)
	EventServiceRestartStart = EventType(iota)
	EventServiceRestart      = EventType(iota)
	EventServicePullStart    = EventType(iota)
	EventServicePull         = EventType(iota)
	EventServiceKillStart    = EventType(iota)
	EventServiceKill         = EventType(iota)
	EventServiceStartStart   = EventType(iota)
	EventServiceStart        = EventType(iota)
	EventServiceBuildStart   = EventType(iota)
	EventServiceBuild        = EventType(iota)
	EventServicePauseStart   = EventType(iota)
	EventServicePause        = EventType(iota)
	EventServiceUnpauseStart = EventType(iota)
	EventServiceUnpause      = EventType(iota)
	EventServiceStopStart    = EventType(iota)
	EventServiceStop         = EventType(iota)
	EventServiceRunStart     = EventType(iota)
	EventServiceRun          = EventType(iota)

	EventProjectDownStart     = EventType(iota)
	EventProjectDownDone      = EventType(iota)
	EventProjectCreateStart   = EventType(iota)
	EventProjectCreateDone    = EventType(iota)
	EventProjectUpStart       = EventType(iota)
	EventProjectUpDone        = EventType(iota)
	EventProjectDeleteStart   = EventType(iota)
	EventProjectDeleteDone    = EventType(iota)
	EventProjectRestartStart  = EventType(iota)
	EventProjectRestartDone   = EventType(iota)
	EventProjectReload        = EventType(iota)
	EventProjectReloadTrigger = EventType(iota)
	EventProjectKillStart     = EventType(iota)
	EventProjectKillDone      = EventType(iota)
	EventProjectStartStart    = EventType(iota)
	EventProjectStartDone     = EventType(iota)
	EventProjectBuildStart    = EventType(iota)
	EventProjectBuildDone     = EventType(iota)
	EventProjectPauseStart    = EventType(iota)
	EventProjectPauseDone     = EventType(iota)
	EventProjectUnpauseStart  = EventType(iota)
	EventProjectUnpauseDone   = EventType(iota)
	EventProjectStopStart     = EventType(iota)
	EventProjectStopDone      = EventType(iota)
)

func (e EventType) String() string {
	var m string
	switch e {
	case EventContainerCreated:
		m = "Created container"
	case EventContainerStarted:
		m = "Started container"

	case EventServiceAdd:
		m = "Adding"
	case EventServiceUpStart:
		m = "Starting"
	case EventServiceUpIgnored:
		m = "Ignoring"
	case EventServiceUp:
		m = "Started"
	case EventServiceCreateStart:
		m = "Creating"
	case EventServiceCreate:
		m = "Created"
	case EventServiceDeleteStart:
		m = "Deleting"
	case EventServiceDelete:
		m = "Deleted"
	case EventServiceStopStart:
		m = "Stopping"
	case EventServiceStop:
		m = "Stopped"
	case EventServiceDownStart:
		m = "Stopping"
	case EventServiceDown:
		m = "Stopped"
	case EventServiceRestartStart:
		m = "Restarting"
	case EventServiceRestart:
		m = "Restarted"
	case EventServicePullStart:
		m = "Pulling"
	case EventServicePull:
		m = "Pulled"
	case EventServiceKillStart:
		m = "Killing"
	case EventServiceKill:
		m = "Killed"
	case EventServiceStartStart:
		m = "Starting"
	case EventServiceStart:
		m = "Started"
	case EventServiceBuildStart:
		m = "Building"
	case EventServiceBuild:
		m = "Built"
	case EventServiceRunStart:
		m = "Executing"
	case EventServiceRun:
		m = "Executed"

	case EventProjectDownStart:
		m = "Stopping project"
	case EventProjectDownDone:
		m = "Project stopped"
	case EventProjectStopStart:
		m = "Stopping project"
	case EventProjectStopDone:
		m = "Project stopped"
	case EventProjectCreateStart:
		m = "Creating project"
	case EventProjectCreateDone:
		m = "Project created"
	case EventProjectUpStart:
		m = "Starting project"
	case EventProjectUpDone:
		m = "Project started"
	case EventProjectDeleteStart:
		m = "Deleting project"
	case EventProjectDeleteDone:
		m = "Project deleted"
	case EventProjectRestartStart:
		m = "Restarting project"
	case EventProjectRestartDone:
		m = "Project restarted"
	case EventProjectReload:
		m = "Reloading project"
	case EventProjectReloadTrigger:
		m = "Triggering project reload"
	case EventProjectKillStart:
		m = "Killing project"
	case EventProjectKillDone:
		m = "Project killed"
	case EventProjectStartStart:
		m = "Starting project"
	case EventProjectStartDone:
		m = "Project started"
	case EventProjectBuildStart:
		m = "Building project"
	case EventProjectBuildDone:
		m = "Project built"
	}

	if m == "" {
		m = fmt.Sprintf("EventType: %d", int(e))
	}

	return m
}

// InfoPart holds key/value strings.
type InfoPart struct {
	Key, Value string
}

// InfoSet holds a list of Info.
type InfoSet []Info

// Info holds a list of InfoPart.
type Info []InfoPart

// Project holds libcompose project information.
type Project struct {
	Name           string
	Configs        *config.Configs
	Files          []string
	ReloadCallback func() error
	context        *Context
	reload         []string
	upCount        int
	listeners      []chan<- Event
	hasListeners   bool
}

// Service defines what a libcompose service provides.
type Service interface {
	Info(qFlag bool) (InfoSet, error)
	Name() string
	Build() error
	Create() error
	Up() error
	Start() error
	Stop() error
	Down() error
	Delete() error
	Restart() error
	Log() error
	Pull() error
	Kill() error
	Config() *config.ServiceConfig
	DependentServices() []ServiceRelationship
	Containers() ([]Container, error)
	Scale(count int) error
	Pause() error
	Unpause() error
	Run(commandParts []string) (int, error)
}

// Container defines what a libcompose container provides.
type Container interface {
	ID() (string, error)
	Name() string
	Port(port string) (string, error)
	IsRunning() (bool, error)
}

// ServiceFactory is an interface factory to create Service object for the specified
// project, with the specified name and service configuration.
type ServiceFactory interface {
	Create(project *Project, name string, serviceConfig *config.ServiceConfig) (Service, error)
}

// ServiceRelationshipType defines the type of service relationship.
type ServiceRelationshipType string

// RelTypeLink means the services are linked (docker links).
const RelTypeLink = ServiceRelationshipType("")

// RelTypeNetNamespace means the services share the same network namespace.
const RelTypeNetNamespace = ServiceRelationshipType("netns")

// RelTypeIpcNamespace means the service share the same ipc namespace.
const RelTypeIpcNamespace = ServiceRelationshipType("ipc")

// RelTypeVolumesFrom means the services share some volumes.
const RelTypeVolumesFrom = ServiceRelationshipType("volumesFrom")

// ServiceRelationship holds the relationship information between two services.
type ServiceRelationship struct {
	Target, Alias string
	Type          ServiceRelationshipType
	Optional      bool
}

// NewServiceRelationship creates a new Relationship based on the specified alias
// and relationship type.
func NewServiceRelationship(nameAlias string, relType ServiceRelationshipType) ServiceRelationship {
	name, alias := NameAlias(nameAlias)
	return ServiceRelationship{
		Target: name,
		Alias:  alias,
		Type:   relType,
	}
}
