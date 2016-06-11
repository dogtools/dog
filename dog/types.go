package dog

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/xsb/dog/util"
)

// TaskList is an array of Tasks
type TaskList []Task

// Task is a representation of a dogfile task.
type Task struct {
	Name        string `json:"task,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    bool   `json:"duration,omitempty"`
	Run         string `json:"run,omitempty"`
	Path        string `json:"path,omitempty"`
}

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
