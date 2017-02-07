package run

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// runCmd embeds and extends exec.Cmd.
type runCmd struct {
	exec.Cmd
	tmpFile string
}

// runCmdProperties defines how a new runCmd needs to be created.
type runCmdProperties struct {
	runner        string
	fileExtension string
	code          string
	workdir       string
	env           []string
}

// Wait waits until the command finishes running and provides exit information.
//
// This method overrites the Wait method that comes from the embedded exec.Cmd
// type, adding the removal of the temporary file.
func (c *runCmd) Wait() error {
	defer func() {
		_ = os.Remove(c.tmpFile)
	}()

	err := c.Cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

// writeTempFile copies the code in a temporary file that will get passed as an
// argument to the runner (as in `sh <tmpFile>`).
func (c *runCmd) writeTempFile(data string, fileExtension string) error {
	dir, err := ioutil.TempDir("", "dog")
	if err != nil {
		return err
	}

	c.tmpFile = fmt.Sprintf("%s/task%s", dir, fileExtension)

	err = ioutil.WriteFile(c.tmpFile, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

// newCmdRunner creates a cmd type runner of the chosen executor.
func newCmdRunner(p runCmdProperties) (Runner, error) {
	if p.code == "" {
		return nil, errors.New("No code specified to run")
	}

	cmd := runCmd{}

	path, err := exec.LookPath(strings.Fields(p.runner)[0])
	if err != nil {
		return nil, err
	}
	cmd.Path = path
	cmd.Args = append(cmd.Args, strings.Fields(p.runner)...)
	err = cmd.writeTempFile(p.code, p.fileExtension)
	if err != nil {
		return nil, err
	}
	cmd.Args = append(cmd.Args, cmd.tmpFile)
	cmd.Dir = p.workdir
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, p.env...)
	cmd.Stdin = os.Stdin

	return &cmd, nil
}
