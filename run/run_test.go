// +build integration

package run

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestShRunner(t *testing.T) {
	runner, err := NewShRunner(`echo "Hello $RUNNER"`, ".", []string{"RUNNER=sh"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello sh"
	if got := outputString; got != want {
		t.Errorf("Expected '%v' but was '%v'", want, got)
	}
}

func TestBashRunner(t *testing.T) {
	runner, err := NewBashRunner(`echo "Hello $RUNNER"`, ".", []string{"RUNNER=bash"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello bash"
	if got := outputString; got != want {
		t.Errorf("Expected '%v' but was '%v'", want, got)
	}
}

func getOutputString(runner Runner) (outputString string, err error) {
	output := new(bytes.Buffer)
	runOut, runErr, err := GetOutputs(runner)
	if err != nil {
		return outputString, fmt.Errorf("%s: %s", err, runErr)
	}
	go io.Copy(output, runOut)
	err = runner.Start()
	if err != nil {
		return
	}
	err = runner.Wait()
	if err != nil {
		return
	}
	return strings.TrimSpace(output.String()), nil
}
