package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
func ParseDogfile(d []byte) (tm types.TaskMap, err error) {
	var tasksToParse []*task

	err = yaml.Unmarshal(d, &tasksToParse)
	if err != nil {
		return
	}

	tm = make(types.TaskMap, len(tasksToParse))
	for _, t := range tasksToParse {
		if _, ok := tm[t.Name]; ok {
			return tm, fmt.Errorf("Duplicated task name %s", t.Name)
		} else {
			task := &types.Task{
				Name:        t.Name,
				Description: t.Description,
				Time:        t.Time,
				Run:         t.Run,
				Executor:    t.Executor,
			}
			if task.Pre, err = parseStringSlice(t.Pre); err != nil {
				return
			}
			if task.Post, err = parseStringSlice(t.Post); err != nil {
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
	var dogfiles []os.FileInfo
	var d []byte

	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		var match bool
		match, err = regexp.MatchString(validDogfileName, file.Name())
		if err != nil {
			return
		}

		if match {
			dogfiles = append(dogfiles, file)
		}
	}

	for _, dogfile := range dogfiles {
		var fileData []byte
		fileData, err = ioutil.ReadFile(dogfile.Name())
		if err != nil {
			return
		}
		d = append(d, fileData...)
	}

	return ParseDogfile(d)
}
