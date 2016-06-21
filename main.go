package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dogtools/dog/execute"
	"github.com/dogtools/dog/parser"
	"github.com/dogtools/dog/types"
)

func printHelp() {
	// TODO write the Help text
	fmt.Println("Dog Help")
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

func main() {
	switch {

	// dog
	case len(os.Args) == 1:
		tm, err := parser.LoadDogFile()
		if err != nil {
			fmt.Println("Error: No valid Dogfile in current directory")
			fmt.Println("Need help? --> dog help")
			fmt.Println("More info ---> https://github.com/dogtools/dog")
		} else {
			printTasks(tm)
		}
		os.Exit(0)

	// dog help
	case len(os.Args) == 2 && os.Args[1] == "help":
		printHelp()
		os.Exit(0)

	// dog <task>
	case len(os.Args) >= 2 && os.Args[1] != "help":
		taskName := os.Args[1]

		tm, err := parser.LoadDogFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		printTaskFooter := true // TODO: user can specify with a flag
		runner, err := execute.NewRunner(tm, printTaskFooter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		runner.Run(taskName)
	}
}
