package dog

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/xsb/dog/util"
)

// Task is a representation of a dogfile task.
type Task struct {
	Name        string
	Description string
	Duration    bool
	Run         []byte
	Path        string
}

// ToDisk saves the task command to a temp script.
func (t *Task) ToDisk() error {
	t.Path = "/tmp/dog-" +
		util.RandString(32, rand.NewSource(time.Now().UnixNano())) +
		t.Name
	if err := ioutil.WriteFile(t.Path, t.Run, 0644); err != nil {
		return err
	}
	return nil
}
