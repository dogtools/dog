package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dogtools/dog"
)

const version = "v0.4.0"

func main() {
	// parse cli arguments
	a, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if a.help {
		printHelp()
		os.Exit(0)
	}

	if a.version {
		printVersion()
		os.Exit(0)
	}

	// parse dogfile
	var d dog.Dogfile
	if err = d.ParseFromDisk(a.directory); err != nil {
		printNoValidDogfile()
		os.Exit(1)
	}
	dog.DeprecationWarnings(os.Stderr)

	if a.taskName != "" {
		if a.info {
			dog.ProvideExtraInfo = true
		}

		if d.Tasks[a.taskName] != nil {
			if a.workdir != "" {
				d.Tasks[a.taskName].Workdir = a.workdir
			}
			if d.Tasks[a.taskName].Workdir == "" {
				d.Tasks[a.taskName].Workdir = a.directory
			}
		} else {
			fmt.Println("Unknown task name:", a.taskName)
			os.Exit(1)
		}

		// generate task chain
		var tc dog.TaskChain
		if err = tc.Generate(d, a.taskName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// run task chain
		err = tc.Run(os.Stdout, os.Stderr)
		if err != nil {
			os.Exit(2)
		}

	} else {
		printTasks(d)
		os.Exit(0)
	}
}

// print tasks with description
func printTasks(d dog.Dogfile) {
	maxCharSize := 0
	for taskName, task := range d.Tasks {
		if task.Description != "" && len(taskName) > maxCharSize {
			maxCharSize = len(taskName)
		}
	}

	var tasks []string
	for taskName, task := range d.Tasks {
		if task.Description != "" {
			tasks = append(tasks, taskName)
		}
	}
	sort.Strings(tasks)

	for _, taskName := range tasks {
		spaces := strings.Repeat(" ", maxCharSize-len(taskName)+2)
		fmt.Printf("%s%s%s\n", taskName, spaces, d.Tasks[taskName].Description)
	}
}
