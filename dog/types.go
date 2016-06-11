package dog

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/xsb/dog/util"
)

// Task is a representation of a dogfile task
type Task struct {
	Name        string `json:"task,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    bool   `json:"duration,omitempty"`
	Run         string `json:"run,omitempty"`
	Path        string `json:"path,omitempty"`
}

// TaskMap is a map in which the key is a task name and the value is a Task object
type TaskMap map[string]Task

// ToDisk saves the task command to a temp script.
func (t *Task) ToDisk() error {
	t.Path = "/tmp/dog-" +
		util.RandString(32, rand.NewSource(time.Now().UnixNano())) +
		t.Name
	if err := ioutil.WriteFile(t.Path, []byte(t.Run), 0644); err != nil {
		return err
	}
	return nil
}
