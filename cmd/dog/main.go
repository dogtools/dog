package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dogtools/dog"
)

const version = "v0.5.0"

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

	if a.debug {
		fmt.Fprintf(os.Stderr, "[dog-debug] version: %s\n", version)
		fmt.Fprintf(os.Stderr, "[dog-debug] info: %v\n", a.info)
	}

	// parse dogfile
	dtasks, err := dog.ParseFromDisk(a.directory)
	if err != nil {
		printNoValidDogfile()
		os.Exit(1)
	}

	if a.debug {
		fmt.Fprintf(os.Stderr, "[dog-debug] path: %s\n", dtasks.Path)
		fmt.Fprintf(os.Stderr, "[dog-debug] dogfiles: %v\n", dtasks.Files)
	}

	if a.taskName != "" {
		if a.info {
			dog.ProvideExtraInfo = true
		}

		if dtasks.Tasks[a.taskName] == nil {
			fmt.Println("Unknown task name:", a.taskName)
			os.Exit(1)
		}

		// generate task chain
		taskChain, err := dog.NewTaskChain(dtasks, a.taskName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if a.debug {
			var chain string
			for _, t := range taskChain.Tasks {
				chain += fmt.Sprintf("%s ", t.Name)
			}
			fmt.Fprintf(os.Stderr, "[dog-debug] chain: %s\n", chain)
		}

		// run task chain
		err = taskChain.Run(os.Stdout, os.Stderr)
		if err != nil {
			os.Exit(2)
		}

	} else {
		printTasks(dtasks)
		os.Exit(0)
	}
}

// print tasks with description
func printTasks(dtasks dog.Dogtasks) {
	maxCharSize := 0
	for taskName, task := range dtasks.Tasks {
		if task.Description != "" && len(taskName) > maxCharSize {
			maxCharSize = len(taskName)
		}
	}

	var tasks []string
	for taskName, task := range dtasks.Tasks {
		if task.Description != "" {
			tasks = append(tasks, taskName)
		}
	}
	sort.Strings(tasks)

	for _, taskName := range tasks {
		separator := strings.Repeat(" ", maxCharSize-len(taskName)+2)
		fmt.Printf("%s%s%s\n", taskName, separator, dtasks.Tasks[taskName].Description)
		if len(dtasks.Tasks[taskName].Pre) > 0 {
			taskSpace := strings.Repeat(" ", len(taskName))
			fmt.Printf("%s%s  <= %s\n", taskSpace, separator, strings.Join(dtasks.Tasks[taskName].Pre[:], " "))
		}
		if len(dtasks.Tasks[taskName].Post) > 0 {
			taskSpace := strings.Repeat(" ", len(taskName))
			fmt.Printf("%s%s  => %s\n", taskSpace, separator, strings.Join(dtasks.Tasks[taskName].Post[:], " "))
		}
	}
}
