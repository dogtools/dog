package main

// This example shows how to define a Dogfile, parse it, generate the task
// chain and finally run it.

import (
	"fmt"
	"os"

	"github.com/dogtools/dog"
)

func main() {

	// Define two tasks in the Dogfile format using YAML
	dogfile := `
- task: hello-dog
  description: Say Hello
  post: hello-world
  code: echo "Hello Dog!"

- task: hello-world
  description: Say Hello Again
  code: echo "Hello World!"
`

	// Parse Dogfile
	var d dog.Dogfile
	err := d.Parse([]byte(dogfile))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate task chain that starts with 'hello-dog' but include both tasks
	var tc dog.TaskChain
	err = tc.Generate(d, "hello-dog")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run task chain
	err = tc.Run(os.Stdout, os.Stderr)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
