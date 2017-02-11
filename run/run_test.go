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

func TestPythonRunner(t *testing.T) {
	runner, err := NewPythonRunner(`import os
print("Hello %s") % os.environ['RUNNER']`, ".", []string{"RUNNER=python"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello python"
	if got := outputString; got != want {
		t.Errorf("Expected '%v' but was '%v'", want, got)
	}
}

func TestRubyRunner(t *testing.T) {
	runner, err := NewRubyRunner(`puts "Hello #{ENV['RUNNER']}"`, ".", []string{"RUNNER=ruby"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello ruby"
	if got := outputString; got != want {
		t.Errorf("Expected '%v' but was '%v'", want, got)
	}
}

func TestPerlRunner(t *testing.T) {
	runner, err := NewPerlRunner(`use Env; print "Hello $RUNNER"`, ".", []string{"RUNNER=perl"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello perl"
	if got := outputString; got != want {
		t.Errorf("Expected '%v' but was '%v'", want, got)
	}
}

func TestNodejsRunner(t *testing.T) {
	runner, err := NewNodejsRunner(`console.log("Hello " + process.env.RUNNER)`, ".", []string{"RUNNER=nodejs"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello nodejs"
	if got := outputString; got != want {
		t.Errorf("Expected '%v' but was '%v'", want, got)
	}
}

func TestGoRunner(t *testing.T) {
	runner, err := NewGoRunner(`package main
import (
    "fmt"
    "os"
)
func main() {
    runner := os.Getenv("RUNNER")
    fmt.Printf("Hello %s", runner)
}`, ".", []string{"RUNNER=go"})
	if err != nil {
		t.Errorf(err.Error())
	}
	outputString, err := getOutputString(runner)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "Hello go"
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
