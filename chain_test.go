package dog

import (
	"bytes"
	"strings"
	"testing"
)

func TestCycleDetection(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Foo is a task that runs Bar after finishing",
				Post:        []string{"bar"},
				Code:        "echo foo",
			},
			"bar": {
				Name:        "bar",
				Description: "Bar is a task that runs Foo before starting",
				Pre:         []string{"foo"},
				Code:        "echo bar",
			},
		},
	}

	if _, err := NewTaskChain(dogfile, "foo"); err == nil {
		t.Errorf("Failed detecting a cycle in a task chain")
	}
}

func TestRunTaskChain(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Foo says 'foo'",
				Runner:      "sh",
				Code:        "echo foo",
			},
		},
	}

	taskChain, err := NewTaskChain(dogfile, "foo")
	if err != nil {
		t.Errorf("Failed generating a task chain: %s", err)
	}

	runOut, runErr := new(bytes.Buffer), new(bytes.Buffer)
	if err = taskChain.Run(runOut, runErr); err != nil {
		t.Errorf("Failed running a task chain: %s", err)
	}

	if got, want := strings.TrimSpace(runOut.String()), "foo"; got != want {
		t.Errorf("Expected %s but was %s", want, got)
	}
}

func TestRunTaskChainNoRunner(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Task without runner",
				Code:        "echo foo",
			},
		},
	}

	taskChain, err := NewTaskChain(dogfile, "foo")
	if err != nil {
		t.Errorf("Failed generating a task chain: %s", err)
	}

	if err = taskChain.Run(new(bytes.Buffer), new(bytes.Buffer)); err == nil {
		t.Errorf("Failed to detect a task without runner")
	}
}

func TestRunTaskChainUnsupportedRunner(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Task using an unknown runner",
				Runner:      "blade",
				Code:        "echo foo",
			},
		},
	}

	taskChain, err := NewTaskChain(dogfile, "foo")
	if err != nil {
		t.Errorf("Failed generating a task chain: %s", err)
	}

	if err = taskChain.Run(new(bytes.Buffer), new(bytes.Buffer)); err == nil {
		t.Errorf("Failed to detect an unsupported runner: %s", err)
	}
}
