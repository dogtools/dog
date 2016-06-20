package main

import (
	"fmt"
	"os"

	"github.com/dogtools/dog/execute"
	"github.com/dogtools/dog/parser"
	"github.com/dogtools/dog/types"
)

func printHelp() {
	// TODO write the Help text
	fmt.Println("Dog Help")
}

func printTasks(tm types.TaskMap) {
	for k, t := range tm {
		fmt.Printf("%s\t%s\n", k, t.Description)
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

		runner, err := execute.NewRunner(tm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		runner.Run(taskName)
	}
}
