package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dogtools/dog/execute"
	"github.com/dogtools/dog/parser"
	"github.com/joho/godotenv"
)

const version = "v0.1.0"

func main() {
	// if .env file exists (in same dir as Dogfile), load values into env
	if _, err := os.Stat(`./.env`); !os.IsNotExist(err) {
		err = godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			os.Exit(1)
		}
	}

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
				pair := strings.SplitN(e, "=", 2)
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

	if a.workdir != "" {
		tm[a.taskName].Workdir = a.workdir
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
