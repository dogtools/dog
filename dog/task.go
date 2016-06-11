package dog

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/ghodss/yaml"

	"github.com/xsb/dog/util"
)

// Task is a representation of a dogfile task.
type Task struct {
	Name        string `json:"task,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    bool   `json:"duration,omitempty"`
	Run         string `json:"run,omitempty"`
	Path        string `json:"path,omitempty"`
}

type taskList []Task

// LoadDogFile finds a Dogfile in disk, parses YAML and returns a map
func LoadDogFile() (tm map[string]Task, err error) {
	var dat []byte
	var tl taskList

	dat, err = ioutil.ReadFile("Dogfile.yml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(dat, &tl)
	if err != nil {
		return
	}

	// TODO create the map while reading the Dogfile
	tm = make(map[string]Task)
	for _, t := range tl {
		tm[t.Name] = t
	}

	return
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
