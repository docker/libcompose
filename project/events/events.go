// Package events holds event structures, methods and functions.
package events

import (
	"time"
)

// Event holds project-wide event informations.
type Event interface {
	String() string
	Service() string
}

// baseEvent contains the base data for all events
type baseEvent struct {
	Event       string `json:"event"`
	ServiceName string `json:"service"`
}

// String returns a string representation of the event
func (b *baseEvent) String() string {
	return b.Event
}

// Service returns the service name for the event
func (b *baseEvent) Service() string {
	return b.ServiceName
}

// NewEvent creates a new Event which matches the Event interface
func NewEvent(event, service string) Event {
	return &baseEvent{
		Event:       event,
		ServiceName: service,
	}
}

// EventFactory creates a new Event
type EventFactory func(service string) Event

// ErrorEventFactory creates a new Event for a specified service and error
type ErrorEventFactory func(service string, err error) Event

// EventWrapper provides a wrapper around EventFactories to allow
// state dependent Event generation
type EventWrapper interface {
	Started(string) Event
	Failed(string, error) Event
	Done(string) Event
	Action() string
}

type eventWrapper struct {
	startedFactory EventFactory
	failedFactory  ErrorEventFactory
	doneFactory    EventFactory
	action         string
}

// Started creates a new event using the provided EventFactory for
// the 'started' condition
func (wrapper *eventWrapper) Started(serviceName string) Event {
	if wrapper.startedFactory != nil {
		return wrapper.startedFactory(serviceName)
	}
	return nil
}

// Failed creates a new event using the provided eventFactory for
// the 'failed' condition
func (wrapper *eventWrapper) Failed(serviceName string, err error) Event {
	if wrapper.failedFactory != nil {
		return wrapper.failedFactory(serviceName, err)
	}
	return nil
}

// Done creates a new event using the provided eventFactory for
// the 'done' condition
func (wrapper *eventWrapper) Done(serviceName string) Event {
	if wrapper.doneFactory != nil {
		return wrapper.doneFactory(serviceName)
	}
	return nil
}

// Action returns the name of the action this wrapper is supporting
func (wrapper *eventWrapper) Action() string {
	return wrapper.action
}

// NewEventWrapper builds a wrapper around the provided EventFactories
func NewEventWrapper(action string, started EventFactory, done EventFactory, failed ErrorEventFactory) EventWrapper {
	return &eventWrapper{
		startedFactory: started,
		failedFactory:  failed,
		doneFactory:    done,
		action:         action,
	}
}

// NewDummyEventWrapper returns an event wrapper which returns nil events
func NewDummyEventWrapper(action string) EventWrapper {
	return &dummyEventWrapper{
		action: action,
	}
}

type dummyEventWrapper struct {
	action string
}

func (*dummyEventWrapper) Started(string) Event {
	return nil
}

func (*dummyEventWrapper) Done(string) Event {
	return nil
}

func (*dummyEventWrapper) Failed(string, error) Event {
	return nil
}

func (w *dummyEventWrapper) Action() string {
	return w.action
}

// Notifier defines the methods an event notifier should have.
type Notifier interface {
	Notify(event Event)
}

// Emitter defines the methods an event emitter should have.
type Emitter interface {
	AddListener(c chan<- Event)
}

// ContainerEvent holds attributes of container events.
type ContainerEvent struct {
	Event
	ID         string            `json:"id"`
	Time       time.Time         `json:"time"`
	Attributes map[string]string `json:"attributes"`
	Type       string            `json:"type"`
}

// ServiceAdd represents a service being added to a project
type ServiceAdd struct {
	Event
}

// NewServiceAddEvent creates a new Service Add event
func NewServiceAddEvent(serviceName string) Event {
	return &ServiceAdd{
		NewEvent("Service Added", serviceName),
	}
}

// VolumeAdd Represents a volume being added to a project service
type VolumeAdd struct {
	Event
	Driver string
}

// NewVolumeAddEvent creates a new Volume Add event
func NewVolumeAddEvent(serviceName, volumeDriver string) Event {
	return &VolumeAdd{
		Event:  NewEvent("Volume Added", serviceName),
		Driver: volumeDriver,
	}
}

// NetworkAdd Represents a network being added to a project service
type NetworkAdd struct {
	Event
	Driver string
}

// NewNetworkAddEvent creates a new Network Add event
func NewNetworkAddEvent(serviceName, networkDriver string) Event {
	return &NetworkAdd{
		Event:  NewEvent("Network Added", serviceName),
		Driver: networkDriver,
	}
}

// ServiceBuildStart Represents a service build starting
type ServiceBuildStart struct {
	Event
}

// NewServiceBuildStartEvent creates a new service build starting event
func NewServiceBuildStartEvent(serviceName string) Event {
	return &ServiceBuildStart{
		NewEvent("Building Service", serviceName),
	}
}

// ServiceBuildDone represents a service build completing
type ServiceBuildDone struct {
	Event
}

// NewServiceBuildDoneEvent creates a new service build done event
func NewServiceBuildDoneEvent(serviceName string) Event {
	return &ServiceBuildDone{
		NewEvent("Service Built", serviceName),
	}
}

// ServiceBuildFailed represents a service build failing
type ServiceBuildFailed struct {
	Event
	err error
}

// NewServiceBuildFailedEvent creates a new service build failed event
func NewServiceBuildFailedEvent(serviceName string, err error) Event {
	return &ServiceBuildFailed{
		Event: NewEvent("Service Built Failed", serviceName),
		err:   err,
	}
}

// ServiceCreateStart represents a service create starting
type ServiceCreateStart struct {
	Event
}

// NewServiceCreateStartEvent creates a new service create starting event
func NewServiceCreateStartEvent(serviceName string) Event {
	return &ServiceCreateStart{
		NewEvent("Creating Service", serviceName),
	}
}

// ServiceCreateDone represents a service create completing
type ServiceCreateDone struct {
	Event
}

// NewServiceCreateDoneEvent creates a new service create done event
func NewServiceCreateDoneEvent(serviceName string) Event {
	return &ServiceCreateDone{
		NewEvent("Service Created", serviceName),
	}
}

// ServiceCreateFailed represents a service create failing
type ServiceCreateFailed struct {
	Event
	err error
}

// NewServiceCreateFailedEvent creates a new service create failed event
func NewServiceCreateFailedEvent(serviceName string, err error) Event {
	return &ServiceCreateFailed{
		Event: NewEvent("Service Create Failed", serviceName),
		err:   err,
	}
}

// ServiceStopStart represents a service stop starting
type ServiceStopStart struct {
	Event
}

// NewServiceStopStartEvent creates a new service stop starting event
func NewServiceStopStartEvent(serviceName string) Event {
	return &ServiceStopStart{
		NewEvent("Stopping Service", serviceName),
	}
}

// ServiceStopDone represents a service stop completing
type ServiceStopDone struct {
	Event
}

// NewServiceStopDoneEvent creates a new service stop done event
func NewServiceStopDoneEvent(serviceName string) Event {
	return &ServiceStopDone{
		Event: NewEvent("Service Stopped", serviceName),
	}
}

// ServiceStopFailed represents a service stop failing
type ServiceStopFailed struct {
	Event
	err error
}

// NewServiceStopFailedEvent creates a new service stop failed event
func NewServiceStopFailedEvent(serviceName string, err error) Event {
	return &ServiceStopFailed{
		Event: NewEvent("Service Stop Failed", serviceName),
		err:   err,
	}
}

// ServiceRestartStart represents a service restart starting
type ServiceRestartStart struct {
	Event
}

// NewServiceRestartStartEvent creates a new service restart starting event
func NewServiceRestartStartEvent(serviceName string) Event {
	return &ServiceRestartStart{
		NewEvent("Restarting Service", serviceName),
	}
}

// ServiceRestartDone represents a service restart completing
type ServiceRestartDone struct {
	Event
}

// NewServiceRestartDoneEvent creates a new service restart done event
func NewServiceRestartDoneEvent(serviceName string) Event {
	return &ServiceRestartDone{
		NewEvent("Service Restarted", serviceName),
	}
}

// ServiceRestartFailed represents a service restart failing
type ServiceRestartFailed struct {
	Event
	err error
}

// NewServiceRestartFailedEvent creates a new service restart failed event
func NewServiceRestartFailedEvent(serviceName string, err error) Event {
	return &ServiceRestartFailed{
		Event: NewEvent("Service Restart Failed", serviceName),
		err:   err,
	}
}

// ServiceStartStart represents a service start starting
type ServiceStartStart struct {
	Event
}

// NewServiceStartStartEvent creates a new service start starting event
func NewServiceStartStartEvent(serviceName string) Event {
	return &ServiceStartStart{
		NewEvent("Starting Service", serviceName),
	}
}

// ServiceStartDone represents a service start completing
type ServiceStartDone struct {
	Event
}

// NewServiceStartDoneEvent creates a new service start done event
func NewServiceStartDoneEvent(serviceName string) Event {
	return &ServiceStartDone{
		NewEvent("Service Started", serviceName),
	}
}

// ServiceStartFailed represents a service start failing
type ServiceStartFailed struct {
	Event
	err error
}

// NewServiceStartFailedEvent creates a new service start failed event
func NewServiceStartFailedEvent(serviceName string, err error) Event {
	return &ServiceStartFailed{
		Event: NewEvent("Service Start Failed", serviceName),
		err:   err,
	}
}

// ServiceRunStart represents a service run starting
type ServiceRunStart struct {
	Event
}

// NewServiceRunStartEvent creates a new service run starting event
func NewServiceRunStartEvent(serviceName string) Event {
	return &ServiceRunStart{
		NewEvent("Running Service", serviceName),
	}
}

// ServiceRunDone represents a service run completing
type ServiceRunDone struct {
	Event
}

// NewServiceRunDoneEvent creates a new service run done event
func NewServiceRunDoneEvent(serviceName string) Event {
	return &ServiceRunDone{
		NewEvent("Service Run", serviceName),
	}
}

// ServiceRunFailed represents a service run failing
type ServiceRunFailed struct {
	Event
	err error
}

// NewServiceRunFailedEvent creates a new service run failed event
func NewServiceRunFailedEvent(serviceName string, err error) Event {
	return &ServiceRunFailed{
		Event: NewEvent("Service Run Failed", serviceName),
		err:   err,
	}
}

// ServiceUpStart represents a service up starting
type ServiceUpStart struct {
	Event
}

// NewServiceUpStartEvent creates a new service up starting event
func NewServiceUpStartEvent(serviceName string) Event {
	return &ServiceUpStart{
		NewEvent("Starting Service", serviceName),
	}
}

// ServiceUpDone represents a service up completing
type ServiceUpDone struct {
	Event
}

// NewServiceUpDoneEvent creates a new service up done event
func NewServiceUpDoneEvent(serviceName string) Event {
	return &ServiceUpDone{
		NewEvent("Service Started", serviceName),
	}
}

// ServiceUpFailed represents a service up failing
type ServiceUpFailed struct {
	Event
	err error
}

// NewServiceUpFailedEvent creates a new service up failed event
func NewServiceUpFailedEvent(serviceName string, err error) Event {
	return &ServiceUpFailed{
		Event: NewEvent("Service Start Failed", serviceName),
		err:   err,
	}
}

// ServiceUpIgnored represents a service up being ignored
type ServiceUpIgnored struct {
	Event
}

// NewServiceUpIgnoredEvent creates a new service up ignore event
func NewServiceUpIgnoredEvent(serviceName string) Event {
	return &ServiceUpIgnored{
		NewEvent("Service Start Ignored", serviceName),
	}
}

// ServicePullStart represents a service pull starting
type ServicePullStart struct {
	Event
}

// NewServicePullStartEvent creates a new service pull starting event
func NewServicePullStartEvent(serviceName string) Event {
	return &ServicePullStart{
		NewEvent("Pulling Service", serviceName),
	}
}

// ServicePullDone represents a service pull completing
type ServicePullDone struct {
	Event
}

// NewServicePullDoneEvent creates a new service pull done event
func NewServicePullDoneEvent(serviceName string) Event {
	return &ServicePullDone{
		NewEvent("Service Pulled", serviceName),
	}
}

// ServicePullFailed represents a service pull failing
type ServicePullFailed struct {
	Event
	err error
}

// NewServicePullFailedEvent creates a new service pull failed event
func NewServicePullFailedEvent(serviceName string, err error) Event {
	return &ServicePullFailed{
		Event: NewEvent("Service Pull Failed", serviceName),
		err:   err,
	}
}

// ServiceDeleteStart represents a service delete starting
type ServiceDeleteStart struct {
	Event
}

// NewServiceDeleteStartEvent creates a new service delete starting event
func NewServiceDeleteStartEvent(serviceName string) Event {
	return &ServiceDeleteStart{
		NewEvent("Deleting Service", serviceName),
	}
}

// ServiceDeleteDone represents a service delete completing
type ServiceDeleteDone struct {
	Event
}

// NewServiceDeleteDoneEvent creates a new service delete done event
func NewServiceDeleteDoneEvent(serviceName string) Event {
	return &ServiceDeleteDone{
		NewEvent("Service Deleted", serviceName),
	}
}

// ServiceDeleteFailed represents a service delete failing
type ServiceDeleteFailed struct {
	Event
	err error
}

// NewServiceDeleteFailedEvent creates a new service delete failed event
func NewServiceDeleteFailedEvent(serviceName string, err error) Event {
	return &ServiceDeleteFailed{
		Event: NewEvent("Service Delete Failed", serviceName),
		err:   err,
	}
}

// ServiceKillStart represents a service kill starting
type ServiceKillStart struct {
	Event
}

// NewServiceKillStartEvent creates a new service kill starting event
func NewServiceKillStartEvent(serviceName string) Event {
	return &ServiceKillStart{
		NewEvent("Killing Service", serviceName),
	}
}

// ServiceKillDone represents a service kill completing
type ServiceKillDone struct {
	Event
}

// NewServiceKillDoneEvent creates a new service kill done event
func NewServiceKillDoneEvent(serviceName string) Event {
	return &ServiceKillDone{
		NewEvent("Service Killed", serviceName),
	}
}

// ServiceKillFailed represents a service kill failing
type ServiceKillFailed struct {
	Event
	err error
}

// NewServiceKillFailedEvent creates a new service kill failed event
func NewServiceKillFailedEvent(serviceName string, err error) Event {
	return &ServiceKillFailed{
		Event: NewEvent("Service Kill Failed", serviceName),
		err:   err,
	}
}

// ServicePauseStart represents a service pause starting
type ServicePauseStart struct {
	Event
}

// NewServicePauseStartEvent creates a new service pause starting event
func NewServicePauseStartEvent(serviceName string) Event {
	return &ServicePauseStart{
		NewEvent("Pausing Service", serviceName),
	}
}

// ServicePauseDone represents a service pause completing
type ServicePauseDone struct {
	Event
}

// NewServicePauseDoneEvent creates a new service pause done event
func NewServicePauseDoneEvent(serviceName string) Event {
	return &ServicePauseDone{
		NewEvent("Service Paused", serviceName),
	}
}

// ServicePauseFailed represents a service pause failing
type ServicePauseFailed struct {
	Event
	err error
}

// NewServicePauseFailedEvent creates a new service pause failed event
func NewServicePauseFailedEvent(serviceName string, err error) Event {
	return &ServicePauseFailed{
		Event: NewEvent("Service Pause Failed", serviceName),
		err:   err,
	}
}

// ServiceUnpauseStart represents a service unpause starting
type ServiceUnpauseStart struct {
	Event
}

// NewServiceUnpauseStartEvent creates a new service unpause starting event
func NewServiceUnpauseStartEvent(serviceName string) Event {
	return &ServiceUnpauseStart{
		NewEvent("Unpausing Service", serviceName),
	}
}

// ServiceUnpauseDone represents a service unpause completing
type ServiceUnpauseDone struct {
	Event
}

// NewServiceUnpauseDoneEvent creates a new service unpause done event
func NewServiceUnpauseDoneEvent(serviceName string) Event {
	return &ServiceUnpauseDone{
		NewEvent("Service Unpaused", serviceName),
	}
}

// ServiceUnpauseFailed represents a service unpause failing
type ServiceUnpauseFailed struct {
	Event
	err error
}

// NewServiceUnpauseFailedEvent creates a new service unpause failed event
func NewServiceUnpauseFailedEvent(serviceName string, err error) Event {
	return &ServiceUnpauseFailed{
		Event: NewEvent("Service Unpause Failed", serviceName),
		err:   err,
	}
}

// ServiceDownStart represents a service down starting
type ServiceDownStart struct {
	Event
}

// NewServiceDownStartEvent creates a new service down starting event
func NewServiceDownStartEvent(serviceName string) Event {
	return &ServiceDownStart{
		NewEvent("Stopping Service", serviceName),
	}
}

// ServiceDownDone represents a service down completing
type ServiceDownDone struct {
	Event
}

// NewServiceDownDoneEvent creates a new service down done event
func NewServiceDownDoneEvent(serviceName string) Event {
	return &ServiceDownDone{
		NewEvent("Service Stopped", serviceName),
	}
}

// ServiceDownFailed represents a service down failing
type ServiceDownFailed struct {
	Event
	err error
}

// NewServiceDownFailedEvent creates a new service down failed event
func NewServiceDownFailedEvent(serviceName string, err error) Event {
	return &ServiceDownFailed{
		Event: NewEvent("Service Stop Failed", serviceName),
		err:   err,
	}
}

// ProjectRestartStart represents a project restart starting
type ProjectRestartStart struct {
	Event
}

// NewProjectRestartStartEvent creates a new project restart starting event
func NewProjectRestartStartEvent(serviceName string) Event {
	return &ProjectRestartStart{
		NewEvent("Restarting Project", serviceName),
	}
}

// ProjectRestartDone represents a project restart completing
type ProjectRestartDone struct {
	Event
}

// NewProjectRestartDoneEvent creates a new project restart done event
func NewProjectRestartDoneEvent(serviceName string) Event {
	return &ProjectRestartDone{
		NewEvent("Project Restarted", serviceName),
	}
}

// ProjectRestartFailed represents a project restart failing
type ProjectRestartFailed struct {
	Event
	err error
}

// NewProjectRestartFailedEvent creates a new project restart failed event
func NewProjectRestartFailedEvent(serviceName string, err error) Event {
	return &ProjectRestartFailed{
		Event: NewEvent("Project Restart Failed", serviceName),
		err:   err,
	}
}

// ProjectStartStart represents a project start starting
type ProjectStartStart struct {
	Event
}

// NewProjectStartStartEvent creates a new project start starting event
func NewProjectStartStartEvent(serviceName string) Event {
	return &ProjectStartStart{
		NewEvent("Starting Project", serviceName),
	}
}

// ProjectStartDone represents a project start completing
type ProjectStartDone struct {
	Event
}

// NewProjectStartDoneEvent creates a new project start done event
func NewProjectStartDoneEvent(serviceName string) Event {
	return &ProjectStartDone{
		NewEvent("Project Started", serviceName),
	}
}

// ProjectStartFailed represents a project start failing
type ProjectStartFailed struct {
	Event
	err error
}

// NewProjectStartFailedEvent creates a new project start failed event
func NewProjectStartFailedEvent(serviceName string, err error) Event {
	return &ProjectStartFailed{
		Event: NewEvent("Project Start Failed", serviceName),
		err:   err,
	}
}

// ProjectUpStart represents a project up starting
type ProjectUpStart struct {
	Event
}

// NewProjectUpStartEvent creates a new project up starting event
func NewProjectUpStartEvent(serviceName string) Event {
	return &ProjectUpStart{
		NewEvent("Starting Project", serviceName),
	}
}

// ProjectUpDone represents a project up completing
type ProjectUpDone struct {
	Event
}

// NewProjectUpDoneEvent creates a new project up done event
func NewProjectUpDoneEvent(serviceName string) Event {
	return &ProjectUpDone{
		NewEvent("Project Started", serviceName),
	}
}

// ProjectUpFailed represents a project up failing
type ProjectUpFailed struct {
	Event
	err error
}

// NewProjectUpFailedEvent creates a new project up failed event
func NewProjectUpFailedEvent(serviceName string, err error) Event {
	return &ProjectUpFailed{
		Event: NewEvent("Project Up Failed", serviceName),
		err:   err,
	}
}

// ProjectDownStart represents a project down starting
type ProjectDownStart struct {
	Event
}

// NewProjectDownStartEvent creates a new project down starting event
func NewProjectDownStartEvent(serviceName string) Event {
	return &ProjectDownStart{
		NewEvent("Stopping Project", serviceName),
	}
}

// ProjectDownDone represents a project down completing
type ProjectDownDone struct {
	Event
}

// NewProjectDownDoneEvent creates a new project down done event
func NewProjectDownDoneEvent(serviceName string) Event {
	return &ProjectDownDone{
		NewEvent("Project Stopped", serviceName),
	}
}

// ProjectDownFailed represents a project down failing
type ProjectDownFailed struct {
	Event
	err error
}

// NewProjectDownFailedEvent creates a new project down failed event
func NewProjectDownFailedEvent(serviceName string, err error) Event {
	return &ProjectDownFailed{
		Event: NewEvent("Project Down Failed", serviceName),
		err:   err,
	}
}

// ProjectDeleteStart represents a project delete starting
type ProjectDeleteStart struct {
	Event
}

// NewProjectDeleteStartEvent creates a new project delete starting event
func NewProjectDeleteStartEvent(serviceName string) Event {
	return &ProjectDeleteStart{
		NewEvent("Deleting Project", serviceName),
	}
}

// ProjectDeleteDone represents a project delete completing
type ProjectDeleteDone struct {
	Event
}

// NewProjectDeleteDoneEvent creates a new project delete done event
func NewProjectDeleteDoneEvent(serviceName string) Event {
	return &ProjectDeleteDone{
		NewEvent("Project Deleted", serviceName),
	}
}

// ProjectDeleteFailed represents a project delete failing
type ProjectDeleteFailed struct {
	Event
	err error
}

// NewProjectDeleteFailedEvent creates a new project delete failed event
func NewProjectDeleteFailedEvent(serviceName string, err error) Event {
	return &ProjectDeleteFailed{
		Event: NewEvent("Project Delete Failed", serviceName),
		err:   err,
	}
}

// ProjectKillStart represents a project kill starting
type ProjectKillStart struct {
	Event
}

// NewProjectKillStartEvent creates a new project kill starting event
func NewProjectKillStartEvent(serviceName string) Event {
	return &ProjectKillStart{
		NewEvent("Killing Project", serviceName),
	}
}

// ProjectKillDone represents a project kill completing
type ProjectKillDone struct {
	Event
}

// NewProjectKillDoneEvent creates a new project kill done event
func NewProjectKillDoneEvent(serviceName string) Event {
	return &ProjectKillDone{
		NewEvent("Project Killed", serviceName),
	}
}

// ProjectKillFailed represents a project kill failing
type ProjectKillFailed struct {
	Event
	err error
}

// NewProjectKillFailedEvent creates a new project kill failed event
func NewProjectKillFailedEvent(serviceName string, err error) Event {
	return &ProjectKillFailed{
		Event: NewEvent("Project Kill Failed", serviceName),
		err:   err,
	}
}

// ProjectPauseStart represents a project pause starting
type ProjectPauseStart struct {
	Event
}

// NewProjectPauseStartEvent creates a new project pause starting event
func NewProjectPauseStartEvent(serviceName string) Event {
	return &ProjectPauseStart{
		NewEvent("Pausing Project", serviceName),
	}
}

// ProjectPauseDone represents a project pause completing
type ProjectPauseDone struct {
	Event
}

// NewProjectPauseDoneEvent creates a new project pause done event
func NewProjectPauseDoneEvent(serviceName string) Event {
	return &ProjectPauseDone{
		NewEvent("Project Paused", serviceName),
	}
}

// ProjectPauseFailed represents a project pause failing
type ProjectPauseFailed struct {
	Event
	err error
}

// NewProjectPauseFailedEvent creates a new project pause failed event
func NewProjectPauseFailedEvent(serviceName string, err error) Event {
	return &ProjectPauseFailed{
		Event: NewEvent("Project Pause Failed", serviceName),
		err:   err,
	}
}

// ProjectUnpauseStart represents a project unpause starting
type ProjectUnpauseStart struct {
	Event
}

// NewProjectUnpauseStartEvent creates a new project unpause starting event
func NewProjectUnpauseStartEvent(serviceName string) Event {
	return &ProjectUnpauseStart{
		NewEvent("Unpausing Project", serviceName),
	}
}

// ProjectUnpauseDone represents a project unpause completing
type ProjectUnpauseDone struct {
	Event
}

// NewProjectUnpauseDoneEvent creates a new project unpause done event
func NewProjectUnpauseDoneEvent(serviceName string) Event {
	return &ProjectUnpauseDone{
		NewEvent("Project Unpaused", serviceName),
	}
}

// ProjectUnpauseFailed represents a project unpause failing
type ProjectUnpauseFailed struct {
	Event
	err error
}

// NewProjectUnpauseFailedEvent creates a new project unpause failed event
func NewProjectUnpauseFailedEvent(serviceName string, err error) Event {
	return &ProjectUnpauseFailed{
		Event: NewEvent("Project Unpause Failed", serviceName),
		err:   err,
	}
}

// ProjectBuildStart represents a project build starting
type ProjectBuildStart struct {
	Event
}

// NewProjectBuildStartEvent creates a new project build starting event
func NewProjectBuildStartEvent(serviceName string) Event {
	return &ProjectBuildStart{
		NewEvent("Building Project", serviceName),
	}
}

// ProjectBuildDone represents a project build completing
type ProjectBuildDone struct {
	Event
}

// NewProjectBuildDoneEvent creates a new project build done event
func NewProjectBuildDoneEvent(serviceName string) Event {
	return &ProjectBuildDone{
		NewEvent("Project Built", serviceName),
	}
}

// ProjectBuildFailed represents a project build failing
type ProjectBuildFailed struct {
	Event
	err error
}

// NewProjectBuildFailedEvent creates a new project build failed event
func NewProjectBuildFailedEvent(serviceName string, err error) Event {
	return &ProjectBuildFailed{
		Event: NewEvent("Project Build Failed", serviceName),
		err:   err,
	}
}

// ProjectCreateStart represents a project creating
type ProjectCreateStart struct {
	Event
}

// NewProjectCreateStartEvent creates a new project creating event
func NewProjectCreateStartEvent(serviceName string) Event {
	return &ProjectCreateStart{
		NewEvent("Creating Project", serviceName),
	}
}

// ProjectCreateDone represents a project create completing
type ProjectCreateDone struct {
	Event
}

// NewProjectCreateDoneEvent creates a new project create done event
func NewProjectCreateDoneEvent(serviceName string) Event {
	return &ProjectCreateDone{
		NewEvent("Project Created", serviceName),
	}
}

// ProjectCreateFailed represents a project create failing
type ProjectCreateFailed struct {
	Event
	err error
}

// NewProjectCreateFailedEvent creates a new project create failed event
func NewProjectCreateFailedEvent(serviceName string, err error) Event {
	return &ProjectCreateFailed{
		Event: NewEvent("Project Create Failed", serviceName),
		err:   err,
	}
}

// ProjectStopStart represents a project stop
type ProjectStopStart struct {
	Event
}

// NewProjectStopStartEvent creates a new project stopping event
func NewProjectStopStartEvent(serviceName string) Event {
	return &ProjectStopStart{
		NewEvent("Stopping Project", serviceName),
	}
}

// ProjectStopDone represents a project stop completing
type ProjectStopDone struct {
	Event
}

// NewProjectStopDoneEvent creates a new project stop done event
func NewProjectStopDoneEvent(serviceName string) Event {
	return &ProjectStopDone{
		NewEvent("Project Stopped", serviceName),
	}
}

// ProjectStopFailed represents a project stop failing
type ProjectStopFailed struct {
	Event
	err error
}

// NewProjectStopFailedEvent creates a new project stop failed event
func NewProjectStopFailedEvent(serviceName string, err error) Event {
	return &ProjectStopFailed{
		Event: NewEvent("Project Stop Failed", serviceName),
		err:   err,
	}
}

// ProjectReloadDone represents a project reload completing
type ProjectReloadDone struct {
	Event
}

// NewProjectReloadDoneEvent creates a new project reload done event
func NewProjectReloadDoneEvent(serviceName string) Event {
	return &ProjectReloadDone{
		NewEvent("Project Reloaded", serviceName),
	}
}

// ProjectReloadTriggered represents a project reload triggered
type ProjectReloadTriggered struct {
	Event
}

// NewProjectReloadTriggeredEvent creates a new project reload triggered event
func NewProjectReloadTriggeredEvent(serviceName string) Event {
	return &ProjectReloadTriggered{
		NewEvent("Reloading Project", serviceName),
	}
}

// ContainerCreateStart represents a container create
type ContainerCreateStart struct {
	Event
	ContainerName string
}

// NewContainerCreateStartEvent creates a new container creating event
func NewContainerCreateStartEvent(serviceName, containerName string) Event {
	return &ContainerCreateStart{
		Event:         NewEvent("Creating Container", serviceName),
		ContainerName: containerName,
	}
}

// ContainerCreateDone represents a container create completing
type ContainerCreateDone struct {
	Event
	ContainerName string
}

// NewContainerCreateDoneEvent creates a new container create done event
func NewContainerCreateDoneEvent(serviceName, containerName string) Event {
	return &ContainerCreateDone{
		Event:         NewEvent("Container Created", serviceName),
		ContainerName: containerName,
	}
}

// ContainerCreateFailed represents a container create failing
type ContainerCreateFailed struct {
	Event
	err           error
	ContainerName string
}

// NewContainerCreateFailedEvent creates a new container create failed event
func NewContainerCreateFailedEvent(serviceName, containerName string, err error) Event {
	return &ContainerCreateFailed{
		Event:         NewEvent("Container Create Failed", serviceName),
		err:           err,
		ContainerName: containerName,
	}
}

// ContainerStartStart represents a container start
type ContainerStartStart struct {
	Event
	ContainerName string
}

// NewContainerStartStartEvent creates a new container starting event
func NewContainerStartStartEvent(serviceName, containerName string) Event {
	return &ContainerStartStart{
		Event:         NewEvent("Container Starting", serviceName),
		ContainerName: containerName,
	}
}

// ContainerStartDone represents a container start completing
type ContainerStartDone struct {
	Event
	ContainerName string
}

// NewContainerStartDoneEvent creates a new container start done event
func NewContainerStartDoneEvent(serviceName, containerName string) Event {
	return &ContainerStartDone{
		Event:         NewEvent("Container Started", serviceName),
		ContainerName: containerName,
	}
}

// ContainerStartFailed represents a container start failing
type ContainerStartFailed struct {
	Event
	err           error
	ContainerName string
}

// NewContainerStartFailedEvent creates a new container start failed event
func NewContainerStartFailedEvent(serviceName, containerName string, err error) Event {
	return &ContainerStartFailed{
		Event:         NewEvent("Container Start Failed", serviceName),
		err:           err,
		ContainerName: containerName,
	}
}
