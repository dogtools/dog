package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/dogtools/dog/types"
	"github.com/ghodss/yaml"
)

func ParseDogfile(d []byte) (tm types.TaskMap, err error) {
	var tasks []*types.Task

	err = yaml.Unmarshal(d, &tasks)
	if err != nil {
		return
	}

	tm = make(types.TaskMap)
	for _, t := range tasks {
		if _, ok := tm[t.Name]; ok {
			return tm, fmt.Errorf("Duplicated task name %s", t.Name)
		} else {
			if pre, ok := t.Pre.(string); ok {
				t.Pre = []string{pre}
			} else if t.Pre == nil {
				t.Pre = []string{}
			} else if _, ok = t.Pre.([]string); !ok {
				return tm, fmt.Errorf("Invalid pre for task %s", t.Name)
			}

			if post, ok := t.Post.(string); ok {
				t.Post = []string{post}
			} else if t.Post == nil {
				t.Post = []string{}
			} else if _, ok = t.Post.([]string); !ok {
				return tm, fmt.Errorf("Invalid post for task %s", t.Name)
			}
			tm[t.Name] = t
		}
	}

	return
}

// LoadDogFile finds a Dogfile in disk, parses YAML and returns a map
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
