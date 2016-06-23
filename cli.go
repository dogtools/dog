package main

import (
	"fmt"
	"strings"

	"github.com/dogtools/dog/types"
)

type userArgs struct {
	help        bool
	version     bool
	printFooter bool
	taskName    string
	taskArgs    []string
}

var knownFlags = [...]string{
	"--help", "-h",
	"--version", "-v",
	"--footer", "--no-footer",
}

func printVersion() {
	fmt.Println("dog version: " + version)
}

func printHelp() {
	fmt.Println(`Usage: dog
       dog [OPTIONS] TASK [ARGS]
       dog [--help] [--version]

Dog is a command line application that executes tasks.

Options:

  --footer       Print information footer after task execution (default)
  --no-footer    Don't print information footer after task execution
  -h, --help     Print this help
  -v, --version  Print version information and quit`)
}

func printNoValidDogfile() {
	fmt.Println(`Error: No valid Dogfile in current directory
Need help? --> dog --help
More info  --> https://github.com/dogtools/dog`)
}

func printTasks(tm types.TaskMap) {
	maxCharSize := 0

	for taskName, _ := range tm {
		if len(taskName) > maxCharSize {
			maxCharSize = len(taskName)
		}
	}

	for taskName, t := range tm {
		spaces := strings.Repeat(" ", maxCharSize-len(taskName)+2)
		fmt.Printf("%s%s%s\n", taskName, spaces, t.Description)
	}
}

func parseArgs(args []string) (a userArgs, err error) {

	// default values
	a = userArgs{
		help:        false,
		version:     false,
		printFooter: true,
		taskName:    "",
		taskArgs:    []string{},
	}

	// iterate over all provided arguments
	for i, arg := range args {

		if arg == "--help" || arg == "-h" {
			if i == 0 && len(args) == 1 && a.taskName == "" {
				a.help = true
				return a, nil
			} else {
				return a, fmt.Errorf("Error: --help doesn't accept additional parameters")
			}
		}

		if arg == "--version" || arg == "-v" {
			if i == 0 && len(args) == 1 && a.taskName == "" {
				a.version = true
				return a, nil
			} else {
				return a, fmt.Errorf("Error: --version doesn't accept additional parameters")
			}
		}

		if arg == "--footer" {
			if a.taskName == "" {
				a.printFooter = true
			} else {
				return a, fmt.Errorf("Error: --footer is not a valid task argument")
			}
		}

		if arg == "--no-footer" {
			if a.taskName == "" {
				a.printFooter = false
			} else {
				return a, fmt.Errorf("Error: --no-footer is not a valid task argument")
			}
		}

		if string(arg[0]) != "-" {
			if a.taskName == "" {
				a.taskName = arg
			} else {
				return a, fmt.Errorf("Error: only one task can be executed at a time")
			}
		} else {
			validArg := false
			for _, f := range knownFlags {
				if arg == f {
					validArg = true
				}
			}
			if !validArg {
				return a, fmt.Errorf("Error: %s is not a valid argument", arg)
			}
		}
	}

	return a, nil
}
