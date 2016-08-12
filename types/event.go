package types

import "time"

type EventType int

const (
	StartEvent = EventType(iota + 1)
	OutputEvent
	EndEvent
)

type Event struct {
	Type     EventType
	Task     string
	Time     time.Time
	Body     []byte
	ExitCode int
}

func NewStartEvent(taskName string) *Event {
	return &Event{
		Type: StartEvent,
		Task: taskName,
		Time: time.Now(),
	}
}

func NewOutputEvent(taskName string, body []byte) *Event {
	return &Event{
		Type: OutputEvent,
		Task: taskName,
		Time: time.Now(),
		Body: body,
	}
}

func NewEndEvent(taskName string, exitCode int) *Event {
	return &Event{
		Type:     EndEvent,
		Task:     taskName,
		Time:     time.Now(),
		ExitCode: exitCode,
	}
}
