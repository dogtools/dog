// Package sh provide an implementation
// for sh command executor.
package sh

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/xsb/dog/dog"
)

// Sh imlements sh shell executor.
type Sh struct {
	cmd string
}

func init() {
	sh := &Sh{}
	sh.cmd = "sh"
	dog.RegisterExecutor("sh", sh)
}

// Exec executes the created tmp script and writes the output to the writer.
func (sh *Sh) Exec(t *dog.Task, w io.Writer) error {
	if err := t.ToDisk(); err != nil {
		return err
	}

	defer func() {
		w.Write([]byte("=== Task " + t.Name + " finished ===\n"))
		// Remove temporary script
		err := os.Remove(t.Path)
		if err != nil {
			panic(err)
		}
	}()

	binary, err := exec.LookPath(sh.cmd)
	if err != nil {
		return err
	}

	cmd := exec.Command(binary, string(t.Run))
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

	err = cmd.Start()
	if err != nil {
		return err
	}
	w.Write([]byte("=== Task " + t.Name + " started ===\n"))

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
