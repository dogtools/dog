// Package def provide an implementation
// for default command shell executor.
package def

import (
	"bufio"
	"io"
	"os"
	"os/exec"

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
	var err error

	if err = t.ToDisk(); err != nil {
		return err
	}

	defer func() {
		// Remove temporary script in goroutine to not block by IO ops.
		go func() {
			err := os.Remove(t.Path)
			if err != nil {
				w.Write([]byte(err.Error() + "\n"))
			}
		}()
	}()

	binary, err := exec.LookPath(def.cmd)
	if err != nil {
		return err
	}

	cmd := exec.Command(binary, t.Path)
	if err = gatherCmdOutput(cmd, w); err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	w.Write([]byte("=== Task " + t.Name + " started ===\n"))

	if err = cmd.Wait(); err != nil {
		return err
	}
	w.Write([]byte("=== Task " + t.Name + " finished ===\n"))

	return nil
}

func gatherCmdOutput(cmd *exec.Cmd, w io.Writer) error {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdoutReader)
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stdoutScanner.Scan() {
			msg := "\033[34m --= MSG: " + stdoutScanner.Text() + "\n\033[0m"
			w.Write([]byte(msg))
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			msg := "\033[31m --= ERR: " + stderrScanner.Text() + "\n\033[0m"
			w.Write([]byte(msg))
		}
	}()

	return nil
}
