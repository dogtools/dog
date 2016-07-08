package types

// Task is a representation of a dogfile task
type Task struct {
	Name        string
	Description string
	Time        bool
	Run         string
	Executor    string
	Pre         []string
	Post        []string
}

// TaskMap is a map in which the key is a task name and the value is a Task object
type TaskMap map[string]*Task
