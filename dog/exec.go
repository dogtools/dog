package dog

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Task struct {
	Name        string
	Description string
	Duration    bool
	Run         []byte
}

func ExecTask(t Task) {
	// TODO check that parameters not empty
	var startTime time.Time

	// Check that executor exists
	binary, err := exec.LookPath("sh")
	if err != nil {
		panic(err)
	}

	// Write the script to disk
	path := "/tmp/dog-" + t.Name + ".sh"
	err = ioutil.WriteFile(path, t.Run, 0644)
	if err != nil {
		panic(err)
	}

	// Define the command to launch and its arguments
	cmd := exec.Command(binary, path)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// Collect and print STDOUT
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	// Start and wait until it finishes
	if t.Duration {
		startTime = time.Now()
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	if t.Duration {
		duration := time.Now().Sub(startTime)
		fmt.Println(duration.Seconds())
	}

	// Remove temporary script
	err = os.Remove(path)
	if err != nil {
		panic(err)
	}

}
