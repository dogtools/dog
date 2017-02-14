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
		t.Fatalf("Failed generating a task chain: %v", err)
	}

	runOut, runErr := new(bytes.Buffer), new(bytes.Buffer)
	if err = taskChain.Run(runOut, runErr); err != nil {
		t.Fatalf("Failed running a task chain: %v", err)
	}

	if got, want := strings.TrimSpace(runOut.String()), "foo"; got != want {
		t.Fatalf("Expected %v but was %v", want, got)
	}
}

func TestRunTaskChainMustFail(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"must-fail": {
				Name:        "must-fail",
				Description: "Returns a non-zero exit status",
				Runner:      "sh",
				Code:        "false",
			},
		},
	}

	taskChain, err := NewTaskChain(dogfile, "must-fail")
	if err != nil {
		t.Fatalf("Failed generating a task chain: %v", err)
	}

	if err = taskChain.Run(new(bytes.Buffer), new(bytes.Buffer)); err == nil {
		t.Fatalf("Failed to detect a non-zero status code")
	}
}

func TestRunTaskChainMultipleHooks(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"pre-task": {
				Name:        "pre-task",
				Description: "This runs before main",
				Runner:      "sh",
				Code:        "echo pre-task",
			},
			"main": {
				Name:        "main",
				Description: "This is the main task",
				Runner:      "sh",
				Pre:         []string{"pre-task"},
				Post:        []string{"post-task", "final-task"},
				Code:        "echo main-task",
			},
			"post-task": {
				Name:        "post-task",
				Description: "This runs after main",
				Runner:      "sh",
				Code:        "echo post-task",
			},
			"final-task": {
				Name:        "final-task",
				Description: "This is the final task",
				Runner:      "sh",
				Code:        "echo final-task",
			},
		},
	}

	taskChain, err := NewTaskChain(dogfile, "main")
	if err != nil {
		t.Fatalf("Failed generating a task chain: %v", err)
	}

	want := []string{"pre-task", "main", "post-task", "final-task"}
	got := []string{}
	for _, t := range taskChain.Tasks {
		got = append(got, t.Name)
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("Expected %v but was %v", want, got)
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
		t.Fatalf("Failed generating a task chain: %v", err)
	}

	if err = taskChain.Run(new(bytes.Buffer), new(bytes.Buffer)); err == nil {
		t.Fatalf("Failed to detect a task without runner")
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
		t.Fatalf("Failed generating a task chain: %v", err)
	}

	if err = taskChain.Run(new(bytes.Buffer), new(bytes.Buffer)); err == nil {
		t.Fatalf("Failed to detect an unsupported runner: %v", err)
	}
}
