package types

// Task is a representation of a dogfile task
type Task struct {
	Name        string      `json:"task"`
	Description string      `json:"description,omitempty"`
	Time        bool        `json:"time,omitempty"`
	Run         string      `json:"run"`
	Executor    string      `json:"exec,omitempty"`
	Pre         interface{} `json:"pre,omitempty"`
	Post        interface{} `json:"post,omitempty"`
}

// TaskMap is a map in which the key is a task name and the value is a Task object
type TaskMap map[string]*Task
