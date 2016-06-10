package main

import (
	"fmt"
	"os"

	"github.com/xsb/dog/dog"
	_ "github.com/xsb/dog/executors"
)

var taskList = map[string]dog.Task{
	"hello": {
		Name:        "hello",
		Description: "Say Hello!",
		Duration:    false,
		Run:         []byte("echo \"hello world\""),
	},
	"bye": {
		Name:        "bye",
		Description: "Good Bye!",
		Duration:    true,
		Run:         []byte("echo bye cruel world"),
	},
	"find": {
		Name:        "find",
		Description: "List all files in $HOME directory",
		Duration:    true,
		Run:         []byte("find /home/xavi"),
	},
}

func main() {
	arg := os.Args[1]
	if arg == "list" || arg == "help" {
		for k, t := range taskList {
			fmt.Printf("%s\t%s\n", k, t.Description)
		}
	} else {
		task := taskList[arg]
		e := dog.GetExecutor("sh")
		e.Exec(&task, os.Stdout)
	}
}
