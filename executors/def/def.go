// Package def provide an implementation
// for default command shell executor.
package def

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/xsb/dog/dog"
)

// Default imlements standard shell executor.
type Default struct {
	cmd string
}

// NewDefaultExecutor returns a default executor with a cmd.
func NewDefaultExecutor(cmd string) *Default {
	return &Default{
		cmd,
	}
}

// Exec executes the created tmp script and writes the output to the writer.
func (def *Default) Exec(t *dog.Task, w io.Writer) error {

	if err := t.ToDisk(); err != nil {
		return err
	}

	binary, err := exec.LookPath(def.cmd)
	if err != nil {
		return err
	}

	cmd := exec.Command(binary, t.Path)

	w.Write([]byte(" - " + t.Name + " started\n"))

	statusCode := 0
	if output, err := cmd.CombinedOutput(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); !ok {
				// For unknown error status codes set it to 1
				statusCode = 1
			} else {
				statusCode = waitStatus.ExitStatus()
			}
		}
		w.Write(output)
		w.Write([]byte("\n" + err.Error() + "\n"))
	} else {
		w.Write(output)
	}

	msg := fmt.Sprintf(" - %s finished with status code %d\n", t.Name, statusCode)
	w.Write([]byte(msg))

	if err := os.Remove(t.Path); err != nil {
		return err
	}

	return nil
}
