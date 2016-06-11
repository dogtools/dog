package dog

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func ParseDogfile(d []byte) (tm TaskMap, err error) {
	var tl TaskList

	err = yaml.Unmarshal(d, &tl)
	if err != nil {
		return
	}

	// TODO create the map while reading the Dogfile
	tm = make(TaskMap)
	for _, t := range tl {
		tm[t.Name] = t
	}

	return
}

// LoadDogFile finds a Dogfile in disk, parses YAML and returns a map
func LoadDogFile() (tm TaskMap, err error) {
	var dat []byte

	dat, err = ioutil.ReadFile("Dogfile.yml")
	if err != nil {
		return
	}

	tm, err = ParseDogfile(dat)
	if err != nil {
		return
	}

	return
}
