package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/dogtools/dog/types"
	"github.com/ghodss/yaml"
)

var ErrMalformedStringArray = errors.New("Malformed strings array")

type task struct {
	Name        string      `json:"task"`
	Description string      `json:"description,omitempty"`
	Time        bool        `json:"time,omitempty"`
	Run         string      `json:"run"`
	Executor    string      `json:"exec,omitempty"`
	Pre         interface{} `json:"pre,omitempty"`
	Post        interface{} `json:"post,omitempty"`
	Env         interface{} `json:"env,omitempty"`
	Workdir     string      `json:"workdir,omitempty"`
}

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

// ParseDogfile takes a byte slice and process it to return a TaskMap.
func ParseDogfile(d []byte, tm types.TaskMap) (err error) {
	const validTaskName = "^[a-z0-9-]+$"
	var tasksToParse []*task

	err = yaml.Unmarshal(d, &tasksToParse)
	if err != nil {
		return
	}

	for _, t := range tasksToParse {
		if _, ok := tm[t.Name]; ok {
			return fmt.Errorf("Duplicated task name %s", t.Name)
		} else if matches, _ := regexp.MatchString(validTaskName, t.Name); !matches {
			return fmt.Errorf("Invalid name for task %s", t.Name)
		} else {
			task := &types.Task{
				Name:        t.Name,
				Description: t.Description,
				Time:        t.Time,
				Run:         t.Run,
				Executor:    t.Executor,
				Workdir:     t.Workdir,
			}
			if task.Pre, err = parseStringSlice(t.Pre); err != nil {
				return
			}
			if task.Post, err = parseStringSlice(t.Post); err != nil {
				return
			}
			if task.Env, err = parseStringSlice(t.Env); err != nil {
				return
			}
			tm[t.Name] = task
		}
	}

	return
}

// LoadDogFile finds a Dogfile in disk, parses YAML and returns a map.
func LoadDogFile() (tm types.TaskMap, err error) {
	const validDogfileName = "^(Dogfile|🐕)"

	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	tm = make(types.TaskMap)
	for _, file := range files {
		var match bool
		match, err = regexp.MatchString(validDogfileName, file.Name())
		if err != nil {
			return
		}

		if match {
			var fileData []byte
			fileData, err = ioutil.ReadFile(file.Name())
			if err != nil {
				return
			}

			if err = ParseDogfile(fileData, tm); err != nil {
				return
			}
		}
	}
	return
}
