package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xsb/dog/dog"
	_ "github.com/xsb/dog/executors"
)

func main() {
	tm, err := dog.LoadDogFile()
	if err != nil {
		log.Fatal(err)
	}

	arg := os.Args[1]
	if arg == "list" || arg == "help" {
		for k, t := range tm {
			fmt.Printf("%s\t%s\n", k, t.Description)
		}
	} else {
		task := tm[arg]
		e := dog.GetExecutor("sh")
		e.Exec(&task, os.Stdout)
	}
}
