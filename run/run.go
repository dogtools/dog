package run

import (
	"bufio"
	"io"
)

// Runner just runs anything.
type Runner interface {
	// StdoutPipe returns a pipe that will be connected to the runner's
	// standard output when the command starts.
	StdoutPipe() (io.ReadCloser, error)

	// StderrPipe returns a pipe that will be connected to the runner's
	// standard error when the command starts.
	StderrPipe() (io.ReadCloser, error)

	// Start starts the runner but does not wait for it to complete.
	Start() error

	// Wait waits for the runner to exit. It must have been started by Start.
	//
	// The returned error is nil if the runner has no problems copying
	// stdin, stdout, and stderr, and exits with a zero exit status.
	Wait() error
}

// NewShRunner creates a system standard shell script runner.
func NewShRunner(code string, workdir string, env []string) (Runner, error) {
	return newCmdRunner("sh", code, workdir, env)
}

// NewBashRunner creates a Bash runner.
func NewBashRunner(code string, workdir string, env []string) (Runner, error) {
	return newCmdRunner("bash", code, workdir, env)
}

// NewPythonRunner creates a Python runner.
func NewPythonRunner(code string, workdir string, env []string) (Runner, error) {
	return newCmdRunner("python", code, workdir, env)
}

// NewRubyRunner creates a Ruby runner.
func NewRubyRunner(code string, workdir string, env []string) (Runner, error) {
	return newCmdRunner("ruby", code, workdir, env)
}

// NewPerlRunner creates a Perl runner.
func NewPerlRunner(code string, workdir string, env []string) (Runner, error) {
	return newCmdRunner("perl", code, workdir, env)
}

// NewNodejsRunner creates a Node.js runner.
func NewNodejsRunner(code string, workdir string, env []string) (Runner, error) {
	return newCmdRunner("node", code, workdir, env)
}

// GetOutputs is a helper method that returns both stdout and stderr outputs
// from the runner.
func GetOutputs(r Runner) (io.Reader, io.Reader, error) {
	stdout, err := r.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := r.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	return bufio.NewReader(stdout), bufio.NewReader(stderr), nil
}
