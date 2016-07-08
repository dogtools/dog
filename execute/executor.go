package execute

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/dogtools/dog/types"
)

// SystemExecutor is the default executor and is platform dependant.
var SystemExecutor *Executor

func init() {
	switch runtime.GOOS {
	case "windows":
		SystemExecutor = NewExecutor("cmd")
	default:
		SystemExecutor = NewExecutor("sh")
	}
}

func writeTempFile(dir, prefix string, data string, perm os.FileMode) (*os.File, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return f, err
	}

	if err = f.Chmod(perm); err != nil {
		return f, err
	}

	_, err = f.WriteString(data)
	return f, err
}

// Executor implements standard shell executor.
type Executor struct {
	cmd string
}

// NewExecutor returns a default executor with a cmd.
func NewExecutor(cmd string) *Executor {
	return &Executor{
		cmd,
	}
}

// Exec executes the created tmp script and writes the output to the writer.
func (ex *Executor) Exec(t *types.Task, eventsChan chan *types.Event) error {
	f, err := writeTempFile("", "dog", t.Run, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			eventsChan <- types.NewOutputEvent(t.Name, []byte(err.Error()))
		}
	}()

	binary, err := exec.LookPath(ex.cmd)
	if err != nil {
		return err
	}

	cmd := exec.Command(binary, f.Name())

	if err := gatherCmdOutput(t.Name, cmd, eventsChan); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return nil
	}

	startEvent := types.NewStartEvent(t.Name)
	eventsChan <- startEvent

	statusCode := 0
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); !ok {
				// For unknown error status codes set it to 1
				statusCode = 1
			} else {
				statusCode = waitStatus.ExitStatus()
			}
		}
	}

	eventsChan <- types.NewEndEvent(t.Name, statusCode, startEvent.Time)

	return nil
}

func gatherCmdOutput(taskName string, cmd *exec.Cmd, eventsChan chan *types.Event) error {
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
			eventsChan <- types.NewOutputEvent(taskName, stdoutScanner.Bytes())
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			eventsChan <- types.NewOutputEvent(taskName, stderrScanner.Bytes())
		}
	}()

	return nil
}
