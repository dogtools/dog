package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xsb/dog/dog"
	_ "github.com/xsb/dog/executors"
)

func printHelp() {
	// TODO write the Help text
	fmt.Println("Dog Help")
}

func printTasks(tm dog.TaskMap) {
	for k, t := range tm {
		fmt.Printf("%s\t%s\n", k, t.Description)
	}
}

func main() {
	switch {

	// dog
	case len(os.Args) == 1:
		tm, err := dog.LoadDogFile()
		if err != nil {
			fmt.Println("Error: No valid Dogfile in current directory")
			fmt.Println("Need help? --> dog help")
			fmt.Println("More info ---> https://github.com/xsb/dog")
		} else {
			printTasks(tm)
		}

	// dog help
	case len(os.Args) == 2 && os.Args[1] == "help":
		printHelp()

	// dog <task>
	case len(os.Args) >= 2 && os.Args[1] != "help":
		taskName := os.Args[1]

		tm, err := dog.LoadDogFile()
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := tm[taskName]; ok {
			task := tm[taskName]
			e := dog.GetExecutor("system")
			e.Exec(&task, os.Stdout)
		} else {
			fmt.Println("No task named " + taskName)
		}

	}
}
