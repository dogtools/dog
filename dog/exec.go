package dog

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func ExecTask(task string, script []byte) (duration time.Duration) {
	// TODO check that parameters not empty

	// Check that executor exists
	binary, err := exec.LookPath("sh")
	if err != nil {
		panic(err)
	}

	// Write the script to disk
	path := "/tmp/dog-" + task + ".sh"
	err = ioutil.WriteFile(path, script, 0644)
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
	startTime := time.Now()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	duration = time.Now().Sub(startTime)

	// Remove temporary script
	err = os.Remove(path)
	if err != nil {
		panic(err)
	}

	return duration
}
