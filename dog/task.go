package dog

import (
	"io/ioutil"
	"os"
)

// Task is a representation of a dogfile task
type Task struct {
	Name        string `json:"task"`
	Description string `json:"description,omitempty"`
	Time        bool   `json:"time,omitempty"`
	Run         string `json:"run"`
	Executor    string `json:"exec"`
	Path        string `json:"-"`
}

// TaskMap is a map in which the key is a task name and the value is a Task object
type TaskMap map[string]Task

func writeTempFile(dir, prefix string, data string, perm os.FileMode) (*os.File, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return f, err
	}

	if err = f.Chmod(perm); err != nil {
		return f, err
	}

	_, err = f.WriteString(data)
	return f, err
}

// ToDisk saves the task command to a temp script.
func (t *Task) ToDisk() (err error) {
	f, err := writeTempFile("", "dog", t.Run, 0644)

	if err == nil {
		t.Path = f.Name()
	}

	return
}
