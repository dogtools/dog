package parser

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/xsb/dog/types"
)

func ParseDogfile(d []byte) (tm types.TaskMap, err error) {
	var tasks []types.Task

	err = yaml.Unmarshal(d, &tasks)
	if err != nil {
		return
	}

	tm = make(types.TaskMap)
	for _, t := range tasks {
		if _, ok := tm[t.Name]; ok {
			// TODO (duplicated task name) fail and return a non-nil error
		} else {
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
