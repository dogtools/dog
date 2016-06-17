package types

import "time"

type Event struct {
	Name   string
	Task   string
	Time   time.Time
	Extras map[string]interface{}
}

func NewStartEvent(taskName string) *Event {
	return &Event{
		Name: "start",
		Task: taskName,
		Time: time.Now(),
	}
}

func NewOutputEvent(taskName string, body []byte) *Event {
	return &Event{
		Name: "output",
		Task: taskName,
		Time: time.Now(),
		Extras: map[string]interface{}{
			"body": body,
		},
	}
}

func NewEndEvent(taskName string, statusCode int) *Event {
	return &Event{
		Name: "end",
		Task: taskName,
		Time: time.Now(),
		Extras: map[string]interface{}{
			"statusCode": statusCode,
		},
	}
}
