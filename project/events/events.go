// Package events holds event structures, methods and functions.
package events

import (
	"time"
)

/*** Event ***/
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

func NewEvent(event, service string) Event {
	return &baseEvent{
		Event:       event,
		ServiceName: service,
	}
}

/*** EventFactory ***/
// EventFactory creates a new Event
type EventFactory func() Event

// EventFactory creates a new Event for a specified service
type ServiceEventFactory func(service string) Event

// ErrorEventFactory creates a new Event for a specified error
type ErrorEventFactory func(err error) Event

// ErrorServiceEventFactory creates a new Event for a specified service and error
type ErrorServiceEventFactory func(service string, err error) Event

/*** EventWrapper ***/
// EventWrapper provides a wrapper around EventFactories to allow
// state dependent Event generation
type EventWrapper interface {
	Started() Event
	Failed(err error) Event
	Done() Event
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
func (wrapper *eventWrapper) Started() Event {
	if wrapper.startedFactory != nil {
		return wrapper.startedFactory()
	} else {
		return nil
	}
}

// Failed creates a new event using the provided EventFactory for
// the 'failed' condition
func (wrapper *eventWrapper) Failed(err error) Event {
	if wrapper.failedFactory != nil {
		return wrapper.failedFactory(err)
	} else {
		return nil
	}
}

// Done creates a new event using the provided EventFactory for
// the 'done' condition
func (wrapper *eventWrapper) Done() Event {
	if wrapper.doneFactory != nil {
		return wrapper.doneFactory()
	} else {
		return nil
	}
}

// Action returns the name of the action this wrapper is supporting
func (wrapper *eventWrapper) Action() string {
	return wrapper.action
}

// NewServiceEventWrapper builds a wrapper around the provided ServiceEventFactories
func NewEventWrapper(action string, started EventFactory, done EventFactory, failed ErrorEventFactory) EventWrapper {
	return &eventWrapper{
		startedFactory: started,
		failedFactory:  failed,
		doneFactory:    done,
		action:         action,
	}
}

// SeviceEventWrapper provides a wrapper around ServiceEventFactories to allow
// state dependent Event generation
type ServiceEventWrapper interface {
	Started(string) Event
	Failed(string, error) Event
	Done(string) Event
	Action() string
}

type serviceEventWrapper struct {
	startedFactory ServiceEventFactory
	failedFactory  ErrorServiceEventFactory
	doneFactory    ServiceEventFactory
	action         string
}

// Started creates a new event using the provided ServiceEventFactory for
// the 'started' condition
func (wrapper *serviceEventWrapper) Started(serviceName string) Event {
	if wrapper.startedFactory != nil {
		return wrapper.startedFactory(serviceName)
	} else {
		return nil
	}
}

// Failed creates a new event using the provided ServiceEventFactory for
// the 'failed' condition
func (wrapper *serviceEventWrapper) Failed(serviceName string, err error) Event {
	if wrapper.failedFactory != nil {
		return wrapper.failedFactory(serviceName, err)
	} else {
		return nil
	}
}

// Done creates a new event using the provided ServiceEventFactory for
// the 'done' condition
func (wrapper *serviceEventWrapper) Done(serviceName string) Event {
	if wrapper.doneFactory != nil {
		return wrapper.doneFactory(serviceName)
	} else {
		return nil
	}
}

// Action returns the name of the action this wrapper is supporting
func (wrapper *serviceEventWrapper) Action() string {
	return wrapper.action
}

// NewServiceEventWrapper builds a wrapper around the provided ServiceEventFactories
func NewServiceEventWrapper(action string, started ServiceEventFactory, done ServiceEventFactory, failed ErrorServiceEventFactory) ServiceEventWrapper {
	return &serviceEventWrapper{
		startedFactory: started,
		failedFactory:  failed,
		doneFactory:    done,
		action:         action,
	}
}

func NewDummyEventWrapper(action string) *dummyEventWrapper {
	return &dummyEventWrapper{
		action: action,
	}
}

type dummyEventWrapper struct {
	action string
}

func (*dummyEventWrapper) Started() Event {
	return nil
}

func (*dummyEventWrapper) Done() Event {
	return nil
}
func (*dummyEventWrapper) Failed(error) Event {
	return nil
}
func (w *dummyEventWrapper) Action() string {
	return w.action
}

func NewDummyServiceEventWrapper(action string) *dummyServiceEventWrapper {
	return &dummyServiceEventWrapper{
		action: action,
	}
}

type dummyServiceEventWrapper struct {
	action string
}

func (*dummyServiceEventWrapper) Started(string) Event {
	return nil
}

func (*dummyServiceEventWrapper) Done(string) Event {
	return nil
}

func (*dummyServiceEventWrapper) Failed(string, error) Event {
	return nil
}

func (w *dummyServiceEventWrapper) Action() string {
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

// Represents a service being added to a project
type ServiceAdd struct {
	*baseEvent
}

// Creates a new Service Add event
func NewServiceAddEvent(serviceName string) Event {
	return &ServiceAdd{
		&baseEvent{
			Event:       "Service Added",
			ServiceName: serviceName,
		},
	}
}

// Represents a volume being added to a project service
type VolumeAdd struct {
	*baseEvent
	Driver string
}

// Creates a new Volume Add event
func NewVolumeAddEvent(serviceName, volumeDriver string) Event {
	return &VolumeAdd{
		baseEvent: &baseEvent{
			Event:       "Volume Added",
			ServiceName: serviceName,
		},
		Driver: volumeDriver,
	}
}

// Represents a network being added to a project service
type NetworkAdd struct {
	*baseEvent
	Driver string
}

// Creates a new Network Add event
func NewNetworkAddEvent(serviceName, networkDriver string) Event {
	return &NetworkAdd{
		baseEvent: &baseEvent{
			Event:       "Network Added",
			ServiceName: serviceName,
		},
		Driver: networkDriver,
	}
}

// Represents a service build starting
type ServiceBuildStart struct {
	*baseEvent
}

// Creates a new service build starting event
func NewServiceBuildStartEvent(serviceName string) Event {
	return &ServiceBuildStart{
		&baseEvent{
			Event:       "Building service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service build completing
type ServiceBuildDone struct {
	*baseEvent
}

// Creates a new service build done event
func NewServiceBuildDoneEvent(serviceName string) Event {
	return &ServiceBuildDone{
		&baseEvent{
			Event:       "Service built",
			ServiceName: serviceName,
		},
	}
}

// Represents a service build failing
type ServiceBuildFailed struct {
	*baseEvent
	err error
}

// Creates a new service build failed event
func NewServiceBuildFailedEvent(serviceName string, err error) Event {
	return &ServiceBuildFailed{
		baseEvent: &baseEvent{
			Event:       "Service build failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service create starting
type ServiceCreateStart struct {
	*baseEvent
}

// Creates a new service create starting event
func NewServiceCreateStartEvent(serviceName string) Event {
	return &ServiceCreateStart{
		&baseEvent{
			Event:       "Creating service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service create completing
type ServiceCreateDone struct {
	*baseEvent
}

// Creates a new service create done event
func NewServiceCreateDoneEvent(serviceName string) Event {
	return &ServiceCreateDone{
		&baseEvent{
			Event:       "Service created",
			ServiceName: serviceName,
		},
	}
}

// Represents a service create failing
type ServiceCreateFailed struct {
	*baseEvent
	err error
}

// Creates a new service create failed event
func NewServiceCreateFailedEvent(serviceName string, err error) Event {
	return &ServiceCreateFailed{
		baseEvent: &baseEvent{
			Event:       "Service create failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service stop starting
type ServiceStopStart struct {
	*baseEvent
}

// Creates a new service stop starting event
func NewServiceStopStartEvent(serviceName string) Event {
	return &ServiceStopStart{
		&baseEvent{
			Event:       "Creating service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service stop completing
type ServiceStopDone struct {
	*baseEvent
}

// Creates a new service stop done event
func NewServiceStopDoneEvent(serviceName string) Event {
	return &ServiceStopDone{
		&baseEvent{
			Event:       "Service stopped",
			ServiceName: serviceName,
		},
	}
}

// Represents a service stop failing
type ServiceStopFailed struct {
	*baseEvent
	err error
}

// Creates a new service stop failed event
func NewServiceStopFailedEvent(serviceName string, err error) Event {
	return &ServiceStopFailed{
		baseEvent: &baseEvent{
			Event:       "Service stop failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service restart starting
type ServiceRestartStart struct {
	*baseEvent
}

// Creates a new service restart starting event
func NewServiceRestartStartEvent(serviceName string) Event {
	return &ServiceRestartStart{
		&baseEvent{
			Event:       "Restarting service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service restart completing
type ServiceRestartDone struct {
	*baseEvent
}

// Creates a new service restart done event
func NewServiceRestartDoneEvent(serviceName string) Event {
	return &ServiceRestartDone{
		&baseEvent{
			Event:       "Service restarted",
			ServiceName: serviceName,
		},
	}
}

// Represents a service restart failing
type ServiceRestartFailed struct {
	*baseEvent
	err error
}

// Creates a new service restart failed event
func NewServiceRestartFailedEvent(serviceName string, err error) Event {
	return &ServiceRestartFailed{
		baseEvent: &baseEvent{
			Event:       "Service restart failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service start starting
type ServiceStartStart struct {
	*baseEvent
}

// Creates a new service start starting event
func NewServiceStartStartEvent(serviceName string) Event {
	return &ServiceStartStart{
		&baseEvent{
			Event:       "Starting service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service start completing
type ServiceStartDone struct {
	*baseEvent
}

// Creates a new service start done event
func NewServiceStartDoneEvent(serviceName string) Event {
	return &ServiceStartDone{
		&baseEvent{
			Event:       "Service started",
			ServiceName: serviceName,
		},
	}
}

// Represents a service start failing
type ServiceStartFailed struct {
	*baseEvent
	err error
}

// Creates a new service start failed event
func NewServiceStartFailedEvent(serviceName string, err error) Event {
	return &ServiceStartFailed{
		baseEvent: &baseEvent{
			Event:       "Service start failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service run starting
type ServiceRunStart struct {
	*baseEvent
}

// Creates a new service run starting event
func NewServiceRunStartEvent(serviceName string) Event {
	return &ServiceRunStart{
		&baseEvent{
			Event:       "Running service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service run completing
type ServiceRunDone struct {
	*baseEvent
}

// Creates a new service run done event
func NewServiceRunDoneEvent(serviceName string) Event {
	return &ServiceRunDone{
		&baseEvent{
			Event:       "Service run",
			ServiceName: serviceName,
		},
	}
}

// Represents a service run failing
type ServiceRunFailed struct {
	*baseEvent
	err error
}

// Creates a new service run failed event
func NewServiceRunFailedEvent(serviceName string, err error) Event {
	return &ServiceRunFailed{
		baseEvent: &baseEvent{
			Event:       "Service run failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service up starting
type ServiceUpStart struct {
	*baseEvent
}

// Creates a new service up starting event
func NewServiceUpStartEvent(serviceName string) Event {
	return &ServiceUpStart{
		&baseEvent{
			Event:       "Starting service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service up completing
type ServiceUpDone struct {
	*baseEvent
}

// Creates a new service up done event
func NewServiceUpDoneEvent(serviceName string) Event {
	return &ServiceUpDone{
		&baseEvent{
			Event:       "Service started",
			ServiceName: serviceName,
		},
	}
}

// Represents a service up failing
type ServiceUpFailed struct {
	*baseEvent
	err error
}

// Creates a new service up failed event
func NewServiceUpFailedEvent(serviceName string, err error) Event {
	return &ServiceUpFailed{
		baseEvent: &baseEvent{
			Event:       "Service start failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service up being ignored
type ServiceUpIgnored struct {
	*baseEvent
}

// Creates a new service up ignore event
func NewServiceUpIgnoredEvent(serviceName string) Event {
	return &ServiceUpFailed{
		baseEvent: &baseEvent{
			Event:       "Service start ignored",
			ServiceName: serviceName,
		},
	}
}

// Represents a service pull starting
type ServicePullStart struct {
	*baseEvent
}

// Creates a new service pull starting event
func NewServicePullStartEvent(serviceName string) Event {
	return &ServicePullStart{
		&baseEvent{
			Event:       "Pulling service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service pull completing
type ServicePullDone struct {
	*baseEvent
}

// Creates a new service pull done event
func NewServicePullDoneEvent(serviceName string) Event {
	return &ServicePullDone{
		&baseEvent{
			Event:       "Service pulled",
			ServiceName: serviceName,
		},
	}
}

// Represents a service pull failing
type ServicePullFailed struct {
	*baseEvent
	err error
}

// Creates a new service pull failed event
func NewServicePullFailedEvent(serviceName string, err error) Event {
	return &ServicePullFailed{
		baseEvent: &baseEvent{
			Event:       "Service pull failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service delete starting
type ServiceDeleteStart struct {
	*baseEvent
}

// Creates a new service delete starting event
func NewServiceDeleteStartEvent(serviceName string) Event {
	return &ServiceDeleteStart{
		&baseEvent{
			Event:       "Deleting service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service delete completing
type ServiceDeleteDone struct {
	*baseEvent
}

// Creates a new service delete done event
func NewServiceDeleteDoneEvent(serviceName string) Event {
	return &ServiceDeleteDone{
		&baseEvent{
			Event:       "Service deleted",
			ServiceName: serviceName,
		},
	}
}

// Represents a service delete failing
type ServiceDeleteFailed struct {
	*baseEvent
	err error
}

// Creates a new service delete failed event
func NewServiceDeleteFailedEvent(serviceName string, err error) Event {
	return &ServiceDeleteFailed{
		baseEvent: &baseEvent{
			Event:       "Service delete failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service kill starting
type ServiceKillStart struct {
	*baseEvent
}

// Creates a new service kill starting event
func NewServiceKillStartEvent(serviceName string) Event {
	return &ServiceKillStart{
		&baseEvent{
			Event:       "Killing service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service kill completing
type ServiceKillDone struct {
	*baseEvent
}

// Creates a new service kill done event
func NewServiceKillDoneEvent(serviceName string) Event {
	return &ServiceKillDone{
		&baseEvent{
			Event:       "Service killed",
			ServiceName: serviceName,
		},
	}
}

// Represents a service kill failing
type ServiceKillFailed struct {
	*baseEvent
	err error
}

// Creates a new service kill failed event
func NewServiceKillFailedEvent(serviceName string, err error) Event {
	return &ServiceKillFailed{
		baseEvent: &baseEvent{
			Event:       "Service kill failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service pause starting
type ServicePauseStart struct {
	*baseEvent
}

// Creates a new service pause starting event
func NewServicePauseStartEvent(serviceName string) Event {
	return &ServicePauseStart{
		&baseEvent{
			Event:       "Pausing service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service pause completing
type ServicePauseDone struct {
	*baseEvent
}

// Creates a new service pause done event
func NewServicePauseDoneEvent(serviceName string) Event {
	return &ServicePauseDone{
		&baseEvent{
			Event:       "Service paused",
			ServiceName: serviceName,
		},
	}
}

// Represents a service pause failing
type ServicePauseFailed struct {
	*baseEvent
	err error
}

// Creates a new service pause failed event
func NewServicePauseFailedEvent(serviceName string, err error) Event {
	return &ServicePauseFailed{
		baseEvent: &baseEvent{
			Event:       "Service pause failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service unpause starting
type ServiceUnpauseStart struct {
	*baseEvent
}

// Creates a new service unpause starting event
func NewServiceUnpauseStartEvent(serviceName string) Event {
	return &ServiceUnpauseStart{
		&baseEvent{
			Event:       "Unpause service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service unpause completing
type ServiceUnpauseDone struct {
	*baseEvent
}

// Creates a new service unpause done event
func NewServiceUnpauseDoneEvent(serviceName string) Event {
	return &ServiceUnpauseDone{
		&baseEvent{
			Event:       "Service unpaused",
			ServiceName: serviceName,
		},
	}
}

// Represents a service unpause failing
type ServiceUnpauseFailed struct {
	*baseEvent
	err error
}

// Creates a new service unpause failed event
func NewServiceUnpauseFailedEvent(serviceName string, err error) Event {
	return &ServiceUnpauseFailed{
		baseEvent: &baseEvent{
			Event:       "Service unpause failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a service down starting
type ServiceDownStart struct {
	*baseEvent
}

// Creates a new service down starting event
func NewServiceDownStartEvent(serviceName string) Event {
	return &ServiceDownStart{
		&baseEvent{
			Event:       "Stopping service",
			ServiceName: serviceName,
		},
	}
}

// Represents a service down completing
type ServiceDownDone struct {
	*baseEvent
}

// Creates a new service down done event
func NewServiceDownDoneEvent(serviceName string) Event {
	return &ServiceDownDone{
		&baseEvent{
			Event:       "Service stopped",
			ServiceName: serviceName,
		},
	}
}

// Represents a service down failing
type ServiceDownFailed struct {
	*baseEvent
	err error
}

// Creates a new service down failed event
func NewServiceDownFailedEvent(serviceName string, err error) Event {
	return &ServiceDownFailed{
		baseEvent: &baseEvent{
			Event:       "Service stop failed",
			ServiceName: serviceName,
		},
		err: err,
	}
}

// Represents a project restart starting
type ProjectRestartStart struct {
	*baseEvent
}

// Creates a new project restart starting event
func NewProjectRestartStartEvent() Event {
	return &ProjectRestartStart{
		&baseEvent{
			Event: "Restarting project",
		},
	}
}

// Represents a project restart completing
type ProjectRestartDone struct {
	*baseEvent
}

// Creates a new project restart done event
func NewProjectRestartDoneEvent() Event {
	return &ProjectRestartDone{
		&baseEvent{
			Event: "Project restarted",
		},
	}
}

// Represents a project restart failing
type ProjectRestartFailed struct {
	*baseEvent
	err error
}

// Creates a new project restart failed event
func NewProjectRestartFailedEvent(err error) Event {
	return &ProjectRestartFailed{
		baseEvent: &baseEvent{
			Event: "Project restart failed",
		},
		err: err,
	}
}

// Represents a project start starting
type ProjectStartStart struct {
	*baseEvent
}

// Creates a new project start starting event
func NewProjectStartStartEvent() Event {
	return &ProjectStartStart{
		&baseEvent{
			Event: "Starting project",
		},
	}
}

// Represents a project start completing
type ProjectStartDone struct {
	*baseEvent
}

// Creates a new project start done event
func NewProjectStartDoneEvent() Event {
	return &ProjectStartDone{
		&baseEvent{
			Event: "Project started",
		},
	}
}

// Represents a project start failing
type ProjectStartFailed struct {
	*baseEvent
	err error
}

// Creates a new project start failed event
func NewProjectStartFailedEvent(err error) Event {
	return &ProjectStartFailed{
		baseEvent: &baseEvent{
			Event: "Project start failed",
		},
		err: err,
	}
}

// Represents a project up starting
type ProjectUpStart struct {
	*baseEvent
}

// Creates a new project up starting event
func NewProjectUpStartEvent() Event {
	return &ProjectUpStart{
		&baseEvent{
			Event: "Starting project",
		},
	}
}

// Represents a project up completing
type ProjectUpDone struct {
	*baseEvent
}

// Creates a new project up done event
func NewProjectUpDoneEvent() Event {
	return &ProjectUpDone{
		&baseEvent{
			Event: "Project started",
		},
	}
}

// Represents a project up failing
type ProjectUpFailed struct {
	*baseEvent
	err error
}

// Creates a new project up failed event
func NewProjectUpFailedEvent(err error) Event {
	return &ProjectUpFailed{
		baseEvent: &baseEvent{
			Event: "Project up failed",
		},
		err: err,
	}
}

// Represents a project down starting
type ProjectDownStart struct {
	*baseEvent
}

// Creates a new project down starting event
func NewProjectDownStartEvent() Event {
	return &ProjectDownStart{
		&baseEvent{
			Event: "Stopping project",
		},
	}
}

// Represents a project down completing
type ProjectDownDone struct {
	*baseEvent
}

// Creates a new project down done event
func NewProjectDownDoneEvent() Event {
	return &ProjectDownDone{
		&baseEvent{
			Event: "Project stopped",
		},
	}
}

// Represents a project down failing
type ProjectDownFailed struct {
	*baseEvent
	err error
}

// Creates a new project down failed event
func NewProjectDownFailedEvent(err error) Event {
	return &ProjectDownFailed{
		baseEvent: &baseEvent{
			Event: "Project down failed",
		},
		err: err,
	}
}

// Represents a project delete starting
type ProjectDeleteStart struct {
	*baseEvent
}

// Creates a new project delete starting event
func NewProjectDeleteStartEvent() Event {
	return &ProjectDeleteStart{
		&baseEvent{
			Event: "Deleting project",
		},
	}
}

// Represents a project delete completing
type ProjectDeleteDone struct {
	*baseEvent
}

// Creates a new project delete done event
func NewProjectDeleteDoneEvent() Event {
	return &ProjectDeleteDone{
		&baseEvent{
			Event: "Project deleted",
		},
	}
}

// Represents a project delete failing
type ProjectDeleteFailed struct {
	*baseEvent
	err error
}

// Creates a new project delete failed event
func NewProjectDeleteFailedEvent(err error) Event {
	return &ProjectDeleteFailed{
		baseEvent: &baseEvent{
			Event: "Project delete failed",
		},
		err: err,
	}
}

// Represents a project kill starting
type ProjectKillStart struct {
	*baseEvent
}

// Creates a new project kill starting event
func NewProjectKillStartEvent() Event {
	return &ProjectKillStart{
		&baseEvent{
			Event: "Killing project",
		},
	}
}

// Represents a project kill completing
type ProjectKillDone struct {
	*baseEvent
}

// Creates a new project kill done event
func NewProjectKillDoneEvent() Event {
	return &ProjectKillDone{
		&baseEvent{
			Event: "Project killed",
		},
	}
}

// Represents a project kill failing
type ProjectKillFailed struct {
	*baseEvent
	err error
}

// Creates a new project kill failed event
func NewProjectKillFailedEvent(err error) Event {
	return &ProjectKillFailed{
		baseEvent: &baseEvent{
			Event: "Project kill failed",
		},
		err: err,
	}
}

// Represents a project pause starting
type ProjectPauseStart struct {
	*baseEvent
}

// Creates a new project pause starting event
func NewProjectPauseStartEvent() Event {
	return &ProjectPauseStart{
		&baseEvent{
			Event: "Pausing project",
		},
	}
}

// Represents a project pause completing
type ProjectPauseDone struct {
	*baseEvent
}

// Creates a new project pause done event
func NewProjectPauseDoneEvent() Event {
	return &ProjectPauseDone{
		&baseEvent{
			Event: "Project paused",
		},
	}
}

// Represents a project pause failing
type ProjectPauseFailed struct {
	*baseEvent
	err error
}

// Creates a new project pause failed event
func NewProjectPauseFailedEvent(err error) Event {
	return &ProjectPauseFailed{
		baseEvent: &baseEvent{
			Event: "Project pause failed",
		},
		err: err,
	}
}

// Represents a project unpause starting
type ProjectUnpauseStart struct {
	*baseEvent
}

// Creates a new project unpause starting event
func NewProjectUnpauseStartEvent() Event {
	return &ProjectUnpauseStart{
		&baseEvent{
			Event: "Unpausing project",
		},
	}
}

// Represents a project unpause completing
type ProjectUnpauseDone struct {
	*baseEvent
}

// Creates a new project unpause done event
func NewProjectUnpauseDoneEvent() Event {
	return &ProjectUnpauseDone{
		&baseEvent{
			Event: "Project unpaused",
		},
	}
}

// Represents a project unpause failing
type ProjectUnpauseFailed struct {
	*baseEvent
	err error
}

// Creates a new project unpause failed event
func NewProjectUnpauseFailedEvent(err error) Event {
	return &ProjectUnpauseFailed{
		baseEvent: &baseEvent{
			Event: "Project unpause failed",
		},
		err: err,
	}
}

// Represents a project build starting
type ProjectBuildStart struct {
	*baseEvent
}

// Creates a new project build starting event
func NewProjectBuildStartEvent() Event {
	return &ProjectBuildStart{
		&baseEvent{
			Event: "Building project",
		},
	}
}

// Represents a project build completing
type ProjectBuildDone struct {
	*baseEvent
}

// Creates a new project build done event
func NewProjectBuildDoneEvent() Event {
	return &ProjectBuildDone{
		&baseEvent{
			Event: "Project built",
		},
	}
}

// Represents a project build failing
type ProjectBuildFailed struct {
	*baseEvent
	err error
}

// Creates a new project build failed event
func NewProjectBuildFailedEvent(err error) Event {
	return &ProjectBuildFailed{
		baseEvent: &baseEvent{
			Event: "Project build failed",
		},
		err: err,
	}
}

// Represents a project creating
type ProjectCreateStart struct {
	*baseEvent
}

// Creates a new project creating event
func NewProjectCreateStartEvent() Event {
	return &ProjectCreateStart{
		&baseEvent{
			Event: "Create project",
		},
	}
}

// Represents a project create completing
type ProjectCreateDone struct {
	*baseEvent
}

// Creates a new project create done event
func NewProjectCreateDoneEvent() Event {
	return &ProjectCreateDone{
		&baseEvent{
			Event: "Project created",
		},
	}
}

// Represents a project create failing
type ProjectCreateFailed struct {
	*baseEvent
	err error
}

// Creates a new project create failed event
func NewProjectCreateFailedEvent(err error) Event {
	return &ProjectCreateFailed{
		baseEvent: &baseEvent{
			Event: "Project create failed",
		},
		err: err,
	}
}

// Represents a project stop
type ProjectStopStart struct {
	*baseEvent
}

// Creates a new project stopping event
func NewProjectStopStartEvent() Event {
	return &ProjectStopStart{
		&baseEvent{
			Event: "Stop project",
		},
	}
}

// Represents a project stop completing
type ProjectStopDone struct {
	*baseEvent
}

// Creates a new project stop done event
func NewProjectStopDoneEvent() Event {
	return &ProjectStopDone{
		&baseEvent{
			Event: "Project stopped",
		},
	}
}

// Represents a project stop failing
type ProjectStopFailed struct {
	*baseEvent
	err error
}

// Creates a new project stop failed event
func NewProjectStopFailedEvent(err error) Event {
	return &ProjectStopFailed{
		baseEvent: &baseEvent{
			Event: "Project stop failed",
		},
		err: err,
	}
}

// Represents a project reload completing
type ProjectReloadDone struct {
	*baseEvent
}

// Creates a new project reload done event
func NewProjectReloadDoneEvent(serviceName string) Event {
	return &ProjectReloadDone{
		&baseEvent{
			Event:       "Project reloaded",
			ServiceName: serviceName,
		},
	}
}

// Represents a project reload triggered
type ProjectReloadTriggered struct {
	*baseEvent
}

// Creates a new project reload triggered event
func NewProjectReloadTriggeredEvent(serviceName string) Event {
	return &ProjectReloadTriggered{
		baseEvent: &baseEvent{
			Event:       "Project reloading",
			ServiceName: serviceName,
		},
	}
}

// Represents a container create
type ContainerCreateStart struct {
	*baseEvent
	ContainerName string
}

// Creates a new container creating event
func NewContainerCreateStartEvent(serviceName, containerName string) Event {
	return &ContainerCreateStart{
		baseEvent: &baseEvent{
			Event:       "Creating container",
			ServiceName: serviceName,
		},
		ContainerName: containerName,
	}
}

// Represents a container create completing
type ContainerCreateDone struct {
	*baseEvent
	ContainerName string
}

// Creates a new container create done event
func NewContainerCreateDoneEvent(serviceName, containerName string) Event {
	return &ContainerCreateDone{
		baseEvent: &baseEvent{
			Event:       "Container created",
			ServiceName: serviceName,
		},
		ContainerName: containerName,
	}
}

// Represents a container create failing
type ContainerCreateFailed struct {
	*baseEvent
	err           error
	ContainerName string
}

// Creates a new container create failed event
func NewContainerCreateFailedEvent(serviceName, containerName string, err error) Event {
	return &ContainerCreateFailed{
		baseEvent: &baseEvent{
			Event:       "Container create failed",
			ServiceName: serviceName,
		},
		err:           err,
		ContainerName: containerName,
	}
}

// Represents a container start
type ContainerStartStart struct {
	*baseEvent
	ContainerName string
}

// Creates a new container starting event
func NewContainerStartStartEvent(serviceName, containerName string) Event {
	return &ContainerStartStart{
		baseEvent: &baseEvent{
			Event:       "Starting container",
			ServiceName: serviceName,
		},
		ContainerName: containerName,
	}
}

// Represents a container start completing
type ContainerStartDone struct {
	*baseEvent
	ContainerName string
}

// Creates a new container start done event
func NewContainerStartDoneEvent(serviceName, containerName string) Event {
	return &ContainerStartDone{
		baseEvent: &baseEvent{
			Event:       "Container started",
			ServiceName: serviceName,
		},
		ContainerName: containerName,
	}
}

// Represents a container start failing
type ContainerStartFailed struct {
	*baseEvent
	err           error
	ContainerName string
}

// Creates a new container start failed event
func NewContainerStartFailedEvent(serviceName, containerName string, err error) Event {
	return &ContainerStartFailed{
		baseEvent: &baseEvent{
			Event:       "Container start failed",
			ServiceName: serviceName,
		},
		err:           err,
		ContainerName: containerName,
	}
}
