package events

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	action      = "foo"
	serviceName = "bar"
)

var (
	testErr = errors.New("test error")
)

func TestNewEvent(t *testing.T) {
	event := NewEvent(action, serviceName)
	assert.Equal(t, serviceName, event.Service())
	assert.Equal(t, action, event.String())
}

func TestEventWrapper(t *testing.T) {
	startedEvent := NewEvent("started", serviceName)
	startedFunc := func() Event {
		return startedEvent
	}
	doneEvent := NewEvent("done", serviceName)
	doneFunc := func() Event {
		return doneEvent
	}
	var failedErr error
	failedEvent := NewEvent("failed", serviceName)
	failedFunc := func(err error) Event {
		failedErr = err
		return failedEvent
	}
	eventWrapper := NewEventWrapper(action, startedFunc, doneFunc, failedFunc)
	assert.Equal(t, action, eventWrapper.Action())

	assert.Equal(t, startedEvent, eventWrapper.Started())
	assert.Equal(t, doneEvent, eventWrapper.Done())
	assert.Equal(t, failedEvent, eventWrapper.Failed(testErr))
	assert.Equal(t, testErr, failedErr)
}

func TestEventWrapperNil(t *testing.T) {
	eventWrapper := NewEventWrapper(action, nil, nil, nil)
	assert.Equal(t, action, eventWrapper.Action())

	assert.Nil(t, eventWrapper.Started())
	assert.Nil(t, eventWrapper.Done())
	assert.Nil(t, eventWrapper.Failed(testErr))
}

func TestServiceEventWrapper(t *testing.T) {
	startedFunc := func(s string) Event {
		return NewEvent("started", s)
	}

	doneFunc := func(s string) Event {
		return NewEvent("done", s)
	}
	var failedErr error
	failedFunc := func(s string, err error) Event {
		failedErr = err
		return NewEvent("failed", s)
	}

	eventWrapper := NewServiceEventWrapper(action, startedFunc, doneFunc, failedFunc)
	assert.Equal(t, action, eventWrapper.Action())

	e1 := eventWrapper.Started(serviceName)
	assert.Equal(t, "started", e1.String())
	assert.Equal(t, serviceName, e1.Service())

	e2 := eventWrapper.Done(serviceName)
	assert.Equal(t, "done", e2.String())
	assert.Equal(t, serviceName, e2.Service())

	e3 := eventWrapper.Failed(serviceName, testErr)
	assert.Equal(t, "failed", e3.String())
	assert.Equal(t, serviceName, e3.Service())
	assert.Equal(t, testErr, failedErr)
}

func TestServiceEventWrapperNil(t *testing.T) {
	eventWrapper := NewServiceEventWrapper(action, nil, nil, nil)
	assert.Equal(t, action, eventWrapper.Action())

	assert.Nil(t, eventWrapper.Started(serviceName))
	assert.Nil(t, eventWrapper.Done(serviceName))
	assert.Nil(t, eventWrapper.Failed(serviceName, testErr))
}

func TestDummyEventWrapper(t *testing.T) {
	eventWrapper := NewDummyEventWrapper(action)
	assert.Equal(t, action, eventWrapper.Action())

	assert.Nil(t, eventWrapper.Started())
	assert.Nil(t, eventWrapper.Done())
	assert.Nil(t, eventWrapper.Failed(testErr))
}

func TestDummyServiceEventWrapper(t *testing.T) {
	eventWrapper := NewDummyServiceEventWrapper(action)
	assert.Equal(t, action, eventWrapper.Action())

	assert.Nil(t, eventWrapper.Started(serviceName))
	assert.Nil(t, eventWrapper.Done(serviceName))
	assert.Nil(t, eventWrapper.Failed(serviceName, testErr))
}
