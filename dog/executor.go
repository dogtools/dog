package dog

import "io"

// Executor is an interface to implement executor commands.
type Executor interface {
	// Exec executes the task and writes command
	// output to the writer.
	Exec(*Task, io.Writer) error
}

// executors holds a relation of the allowed executors.
var executors = map[string]interface{}{}

// RegisterExecutor adds an executor to the registry.
func RegisterExecutor(cmd string, e interface{}) {
	if _, found := executors[cmd]; found {
		panic(ErrAlreadyRegistered)
	}
	executors[cmd] = e
}

// GetExecutor returns an executor initialized with a task.
func GetExecutor(cmd string) Executor {
	if e, found := executors[cmd]; found {
		if ec, ok := e.(Executor); ok {
			return ec
		}
	}
	return nil
}
