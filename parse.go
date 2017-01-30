package dog

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/ghodss/yaml"
)

// ErrMalformedStringArray means that a task have a value of
// pre, post or env that can't be parsed as an array of strings.
var ErrMalformedStringArray = errors.New("Malformed strings array")

// ErrNoDogfile means that the application is unable to find
// a Dogfile in the specified directory.
var ErrNoDogfile = errors.New("No dogfile found")

// Dogfile contains tasks defined in the Dogfile format.
type Dogfile struct {

	// Tasks is used to map task objects by their name.
	Tasks map[string]*Task

	// Path is an optional field that stores the directory
	// where the Dogfile is found.
	Path string

	// Files is an optional field that stores the full path
	// of each Dogfile used to define the Dogfile object.
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

	Pre     interface{} `json:"pre,omitempty"`
	Post    interface{} `json:"post,omitempty"`
	Env     interface{} `json:"env,omitempty"`
	Workdir string      `json:"workdir,omitempty"`
}

// Parse accepts a slice of bytes and parses it following the Dogfile Spec.
func (d *Dogfile) Parse(p []byte) error {
	var tasks []*taskYAML

	err := yaml.Unmarshal(p, &tasks)
	if err != nil {
		return err
	}

	for _, parsedTask := range tasks {
		if _, ok := d.Tasks[parsedTask.Name]; ok {
			return fmt.Errorf("Duplicated task name %s", parsedTask.Name)
		} else if !validTaskName(parsedTask.Name) {
			return fmt.Errorf("Invalid name for task %s", parsedTask.Name)
		} else {
			task := &Task{
				Name:        parsedTask.Name,
				Description: parsedTask.Description,
				Code:        parsedTask.Code,
				Runner:      parsedTask.Runner,
				Workdir:     parsedTask.Workdir,
			}

			// convert pre-tasks, post-tasks and environment variables
			// into []string
			if task.Pre, err = parseStringSlice(parsedTask.Pre); err != nil {
				return err
			}
			if task.Post, err = parseStringSlice(parsedTask.Post); err != nil {
				return err
			}
			if task.Env, err = parseStringSlice(parsedTask.Env); err != nil {
				return err
			}

			// backwards compatibility support for 'run' and 'exec', now called
			// 'code' and 'runner' respectively.
			if parsedTask.Code == "" && parsedTask.Run != "" {
				deprecationWarningRun = true
				task.Code = parsedTask.Run
			}
			if parsedTask.Runner == "" && parsedTask.Exec != "" {
				deprecationWarningExec = true
				task.Runner = parsedTask.Exec
			}

			// set default runner if not specified
			if task.Runner == "" {
				task.Runner = DefaultRunner
			}

			if d.Tasks == nil {
				d.Tasks = make(map[string]*Task)
			}
			d.Tasks[task.Name] = task
		}
	}

	return nil
}

// DeprecationWarnings writes deprecation warnings if they have been found on
// parse time.
//
// Call it with os.Stderr as a parameter to print warnings to STDERR.
func DeprecationWarnings(w io.Writer) {
	if deprecationWarningRun {
		fmt.Fprintln(w,
			"dog: 'run' directive will be deprecated in v0.6.0, use 'code' instead.")
	}
	if deprecationWarningExec {
		fmt.Fprintln(w,
			"dog: 'exec' directive will be deprecated in v0.6.0, use 'runner' instead.")
	}
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
func (d *Dogfile) ParseFromDisk(dir string) error {
	if dir == "" {
		dir = "."
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	d.Path = dir

	files, err := FindDogfiles(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return ErrNoDogfile
	}
	d.Files = files

	for _, file := range d.Files {
		var fileData []byte
		fileData, err = ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		if err = d.Parse(fileData); err != nil {
			return err
		}
	}

	return nil
}

// Validate checks that all tasks in a Dogfile are valid.
func (d *Dogfile) Validate() error {
	for _, t := range d.Tasks {
		if err := t.Validate(); err != nil {
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
	match, err := regexp.MatchString("^(Dogfile|🐕)", name)
	if err != nil {
		return false
	}
	return match
}

// validTaskName checks if a task name is valid as defined
// by the Dogfile Spec.
func validTaskName(name string) bool {
	var match bool
	match, err := regexp.MatchString("^[a-z0-9-]+$", name)
	if err != nil {
		return false
	}
	return match
}
