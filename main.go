package main

import (
	"fmt"
	"os"

	"github.com/xsb/dog/dog"
)

type task struct {
	task        string
	description string
	run         []byte
}

var taskList = map[string]task{
	"hello": task{
		description: "Say Hello!",
		run:         []byte("echo hello world"),
	},
	"bye": task{
		description: "Good Bye!",
		run:         []byte("echo bye cruel world"),
	},
	"find": task{
		description: "List all files in $HOME directory",
		run:         []byte("find /home/xavi"),
	},
}

func main() {
	arg := os.Args[1]
	if arg == "list" || arg == "help" {
		for k, t := range taskList {
			fmt.Printf("%s\t%s\n", k, t.description)
		}
	} else {
		// TODO check that task exists
		task := taskList[arg].task
		run := taskList[arg].run
		duration := dog.ExecTask(task, run)
		fmt.Println(duration.Seconds())
	}
}
