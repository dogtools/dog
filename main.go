package main

import (
	"fmt"
	"os"

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

	tm, err := parser.LoadDogFile()
	if err != nil {
		printNoValidDogfile()
		os.Exit(1)
	}

	if a.taskName != "" {
		runner, err := execute.NewRunner(tm, a.printFooter)
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
