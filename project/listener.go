package project

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project/events"
)

type defaultListener struct {
	project    *Project
	listenChan chan events.Event
	upCount    int
}

// NewDefaultListener create a default listener for the specified project.
func NewDefaultListener(p *Project) chan<- events.Event {
	l := defaultListener{
		listenChan: make(chan events.Event),
		project:    p,
	}
	go l.start()
	return l.listenChan
}

func (d *defaultListener) start() {
	for event := range d.listenChan {
		data, err := json.Marshal(event)

		switch event.(type) {
		case *events.ServiceUpDone:
			d.upCount++
		}

		logf := logrus.Debugf

		if infoLevel(event) {
			logf = logrus.Infof
		}

		if err != nil {
			logf("Failed to Marshal Event [%s] for Project %s. Error: [%s]", event.String(), d.project.Name, err.Error())

		} else if event.Service() == "" {
			logf("Project [%s]: %s %s", d.project.Name, event.String(), data)
		} else {
			logf("[%d/%d] [%s]: %s %s", d.upCount, d.project.ServiceConfigs.Len(), event.Service(), event.String(), data)
		}
	}
}

func infoLevel(event events.Event) bool {
	switch event.(type) {
	case *events.ServiceDeleteStart:
		return true
	case *events.ServiceDeleteDone:
		return true
	case *events.ServiceDeleteFailed:
		return true
	case *events.ServiceDownStart:
		return true
	case *events.ServiceDownDone:
		return true
	case *events.ServiceDownFailed:
		return true
	case *events.ServiceStopStart:
		return true
	case *events.ServiceStopDone:
		return true
	case *events.ServiceStopFailed:
		return true
	case *events.ServiceKillStart:
		return true
	case *events.ServiceKillDone:
		return true
	case *events.ServiceKillFailed:
		return true
	case *events.ServiceCreateStart:
		return true
	case *events.ServiceCreateDone:
		return true
	case *events.ServiceCreateFailed:
		return true
	case *events.ServiceStartStart:
		return true
	case *events.ServiceStartDone:
		return true
	case *events.ServiceStartFailed:
		return true
	case *events.ServiceRestartStart:
		return true
	case *events.ServiceRestartDone:
		return true
	case *events.ServiceRestartFailed:
		return true
	case *events.ServiceUpStart:
		return true
	case *events.ServiceUpDone:
		return true
	case *events.ServiceUpFailed:
		return true
	case *events.ServicePauseStart:
		return true
	case *events.ServicePauseDone:
		return true
	case *events.ServicePauseFailed:
		return true
	case *events.ServiceUnpauseStart:
		return true
	case *events.ServiceUnpauseDone:
		return true
	case *events.ServiceUnpauseFailed:
		return true
	default:
		return false
	}
}
