package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dogtools/dog/execute"
	"github.com/dogtools/dog/parser"
)

const version = "0.0"

func main() {

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

	for a, v := range a.taskArgs {
		switch a {
		case "-e", "--env":
			for _, e := range v {
				pair := strings.Split(e, "=")
				if len(pair) != 2 {
					fmt.Println("Error in env parameter", e)
					os.Exit(1)
				}
				os.Setenv(pair[0], pair[1])
			}
		default:
			fmt.Println("Argument not accepted", a, v)
		}
	}

	tm, err := parser.LoadDogFile()
	if err != nil {
		printNoValidDogfile()
		os.Exit(1)
	}

	if a.taskName != "" {
		runner, err := execute.NewRunner(tm, a.info)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		runner.Run(a.taskName)
	} else {
		printTasks(tm)
		os.Exit(0)
	}
}
