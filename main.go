package main

import (
	"fmt"
	"os"

	"github.com/xsb/dog/executor"
	"github.com/xsb/dog/parser"
	"github.com/xsb/dog/types"
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
			fmt.Println("More info ---> https://github.com/xsb/dog")
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

		if task, ok := tm[taskName]; ok {
			ec := make(chan *types.Event)

			var e *executor.Executor
			if task.Executor != "" {
				e = executor.NewExecutor(task.Executor)
			} else {
				e = executor.SystemExecutor
			}

			go func() {
				if err := e.Exec(&task, ec); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}()

			for {
				select {
				case event := <-ec:
					switch event.Name {
					case "start":
						fmt.Println(" - " + event.Task + " started")
					case "output":
						if body, ok := event.Extras["body"].([]byte); ok {
							fmt.Println(string(body))
						}
					case "end":
						if statusCode, ok := event.Extras["statusCode"].(int); ok {
							fmt.Println(
								fmt.Sprintf(" - %s finished with status code %d", event.Task, statusCode),
							)
							os.Exit(statusCode)
						}
						os.Exit(1)
					}
				}
			}
		} else {
			fmt.Println("No task named " + taskName)
			os.Exit(1)
		}
	}
}
