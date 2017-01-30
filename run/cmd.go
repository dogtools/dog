package run

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

// runCmd embeds and extends exec.Cmd.
type runCmd struct {
	exec.Cmd
	tmpFile *os.File
}

// Wait waits until the command finishes running and provides exit information.
//
// This method overrites the Wait method that comes from the embedded exec.Cmd
// type, adding the removal of the temporary file.
func (c *runCmd) Wait() error {
	defer func() {
		_ = os.Remove(c.tmpFile.Name())
	}()

	err := c.Cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

// writeTempFile copies the code in a temporary file that will get passed as an
// argument to the runner (as in `sh <tmpFile>`).
func (c *runCmd) writeTempFile(data string) error {
	f, err := ioutil.TempFile("", "dog")
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	c.tmpFile = f
	return nil
}

// newCmdRunner creates a cmd type runner of the chosen executor.
func newCmdRunner(runner string, code string, workdir string, env []string) (Runner, error) {
	if code == "" {
		return nil, errors.New("No code specified to run")
	}

	cmd := runCmd{}

	path, err := exec.LookPath(runner)
	if err != nil {
		return nil, err
	}
	cmd.Path = path

	err = cmd.writeTempFile(code)
	if err != nil {
		return nil, err
	}

	cmd.Stdin = os.Stdin
	cmd.Dir = workdir
	cmd.Args = append(cmd.Args, runner, cmd.tmpFile.Name())
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)

	return &cmd, nil
}
