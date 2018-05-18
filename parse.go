package dog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/ghodss/yaml"
)

// DefaultRunner defines the runner to use in case the task does not specify it.
var DefaultRunner = "sh"

// ErrMalformedStringArray means that a task have a value of
// pre, post or env that can't be parsed as an array of strings.
var ErrMalformedStringArray = errors.New("Malformed strings array")

// ErrNoDogfile means that the application is unable to find
// a Dogfile in the specified directory.
var ErrNoDogfile = errors.New("No dogfile found")

// Dogtasks is a collection of tasks with optional metadata from the runtime.
type Dogtasks struct {

	// Tasks is used to map task objects by their name.
	Tasks map[string]*Task

	// Path is an optional field that stores the directory
	// where the Dogfile is found.
	Path string

	// Files is an optional field that stores the full path
	// of each Dogfile used to define the Dogtasks object.
	Files []string
}

// TaskYAML represents a task written in the Dogfile format.
type taskYAML struct {
	Name        string `json:"task"`
	Description string `json:"description,omitempty"`

	Code string `json:"code"`
	Run  string `json:"run"` // backwards compatibility for 'code'

	Runner string `json:"runner,omitempty"`
	Exec   string `json:"exec,omitempty"` // backwards compatibility for 'runner'

	Pre  interface{} `json:"pre,omitempty"`
	Post interface{} `json:"post,omitempty"`
	Env  interface{} `json:"env,omitempty"`

	Workdir  string `json:"workdir,omitempty"`
	Register string `json:"register,omitempty"`
}

// Parse accepts a slice of bytes and parses it following the Dogfile Spec.
func Parse(p []byte) (dtasks Dogtasks, err error) {
	var tasks []*taskYAML

	err = yaml.Unmarshal(p, &tasks)
	if err != nil {
		return
	}

	for _, parsedTask := range tasks {
		if _, ok := dtasks.Tasks[parsedTask.Name]; ok {
			err = fmt.Errorf("Duplicated task name %s", parsedTask.Name)
			return
		} else if !validTaskName(parsedTask.Name) {
			err = fmt.Errorf("Invalid name for task %s", parsedTask.Name)
			return
		} else {
			task := &Task{
				Name:        parsedTask.Name,
				Description: parsedTask.Description,
				Code:        parsedTask.Code,
				Runner:      parsedTask.Runner,
				Workdir:     parsedTask.Workdir,
				Register:    parsedTask.Register,
			}

			// convert pre-tasks, post-tasks and environment variables
			// into []string
			if task.Pre, err = parseStringSlice(parsedTask.Pre); err != nil {
				return
			}
			if task.Post, err = parseStringSlice(parsedTask.Post); err != nil {
				return
			}
			if task.Env, err = parseStringSlice(parsedTask.Env); err != nil {
				return
			}

			// set default runner if not specified
			if task.Runner == "" {
				task.Runner = DefaultRunner
			}

			if dtasks.Tasks == nil {
				dtasks.Tasks = make(map[string]*Task)
			}
			dtasks.Tasks[task.Name] = task
		}
	}

	// validate resulting dogtasks object
	err = dtasks.Validate()

	return
}

// parseStringSlice takes an interface from a pre, post or env field
// and returns a slice of strings representing the found values.
func parseStringSlice(str interface{}) ([]string, error) {
	switch h := str.(type) {
	case string:
		return []string{h}, nil
	case []interface{}:
		s := make([]string, len(h))
		for i, hook := range h {
			sHook, ok := hook.(string)
			if !ok {
				return nil, ErrMalformedStringArray
			}
			s[i] = sHook
		}
		return s, nil
	case nil:
		return []string{}, nil
	default:
		return nil, ErrMalformedStringArray
	}
}

// ParseFromDisk finds a Dogfile in disk and parses it.
func ParseFromDisk(dir string) (dtasks Dogtasks, err error) {
	if dir == "" {
		dir = "."
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	dtasks.Path = dir

	dtasks.Files, err = FindDogfiles(dir)
	if err != nil {
		return
	}
	if len(dtasks.Files) == 0 {
		err = ErrNoDogfile
		return
	}

	// iterate over every found file
	for _, file := range dtasks.Files {
		var fileData []byte
		var d Dogtasks

		fileData, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}

		// parse file
		d, err = Parse(fileData)
		if err != nil {
			return
		}

		// add parsed tasks to main dogfile
		for _, t := range d.Tasks {
			if dtasks.Tasks == nil {
				dtasks.Tasks = make(map[string]*Task)
			}
			dtasks.Tasks[t.Name] = t
		}
	}

	// validate resulting dogfile
	err = dtasks.Validate()

	return
}

// Validate checks that all tasks in a Dogfile are valid.
//
// It checks if any task has a non standard name and also if the
// resulting task chain of each of them have an undesired cycle.
func (dtasks *Dogtasks) Validate() error {
	for _, t := range dtasks.Tasks {

		if !validTaskName(t.Name) {
			return fmt.Errorf("Invalid name for task %s", t.Name)
		}

		if _, err := NewTaskChain(*dtasks, t.Name); err != nil {
			return err
		}

	}
	return nil
}

// FindDogfiles finds Dogfiles in disk for a given path.
//
// It traverses directories until it finds one containing Dogfiles.
// If such a directory is found, the function returns the full path
// for each valid Dogfile in that directory.
func FindDogfiles(p string) ([]string, error) {
	var dogfilePaths []string

	currentPath, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}

	for {
		var files []os.FileInfo
		files, err = ioutil.ReadDir(currentPath)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if validDogfileName(file.Name()) {
				dogfilePath := path.Join(currentPath, file.Name())
				dogfilePaths = append(dogfilePaths, dogfilePath)
			}
		}

		if len(dogfilePaths) > 0 {
			return dogfilePaths, nil
		}

		nextPath := path.Dir(currentPath)
		if nextPath == currentPath {
			return dogfilePaths, nil
		}
		currentPath = nextPath
	}
}

// validDogfileName checks if a Dogfile name is valid as defined
// by the Dogfile Spec.
func validDogfileName(name string) bool {
	var match bool
	match, err := regexp.MatchString("^(dog|🐕).*\\.(yml|yaml)$", name)
	if err != nil {
		return false
	}
	return match
}

// validTaskName checks if a task name is valid as defined
// by the Dogfile Spec.
func validTaskName(name string) bool {
	var match bool
	match, err := regexp.MatchString("^[a-z0-9]+(-[a-z0-9]+)*$", name)
	if err != nil {
		return false
	}
	return match
}
