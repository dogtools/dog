package dog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/dogtools/dog/run"
)

// ProvideExtraInfo specifies if dog needs to provide execution info (duration,
// exit status) after task execution.
var ProvideExtraInfo bool

// ErrCycleInTaskChain means that there is a loop in the path of tasks execution.
var ErrCycleInTaskChain = errors.New("TaskChain includes a cycle of tasks")

// TaskChain contains one or more tasks to be executed in order.
type TaskChain struct {
	Tasks []Task
}

// NewTaskChain creates the task chain for a specific dogfile and task.
func NewTaskChain(dtasks Dogtasks, task string) (taskChain TaskChain, err error) {
	err = taskChain.generate(dtasks, task)
	if err != nil {
		return
	}
	return
}

// Generate recursively iterates over all tasks, including pre and post tasks for
// each of them, and adds all of them into a task chain.
func (taskChain *TaskChain) generate(dtasks Dogtasks, task string) error {

	t, found := dtasks.Tasks[task]
	if !found {
		return fmt.Errorf("Task %q does not exist", task)
	}

	// Cycle detection
	for i := 0; i < len(taskChain.Tasks); i++ {
		if taskChain.Tasks[i].Name == task {
			if len(taskChain.Tasks[i].Pre) > 0 || len(taskChain.Tasks[i].Post) > 0 {
				return ErrCycleInTaskChain
			}
		}
	}

	// Iterate over pre-tasks
	if err := addToChain(taskChain, dtasks, t.Pre); err != nil {
		return err
	}

	// Add current task to chain
	taskChain.Tasks = append(taskChain.Tasks, *t)

	// Iterate over post-tasks
	if err := addToChain(taskChain, dtasks, t.Post); err != nil {
		return err
	}
	return nil
}

// addToChain adds found tasks into the task chain.
func addToChain(taskChain *TaskChain, dtasks Dogtasks, tasks []string) error {
	for _, name := range tasks {

		t, found := dtasks.Tasks[name]
		if !found {
			return fmt.Errorf("Task %q does not exist", name)
		}

		if err := taskChain.generate(dtasks, t.Name); err != nil {
			return err
		}
	}
	return nil
}

// Run handles the execution of all tasks in the TaskChain.
func (taskChain *TaskChain) Run(stdout, stderr io.Writer) error {
	var startTime time.Time
	var registers []string

	for _, t := range taskChain.Tasks {
		var err error
		var runner run.Runner
		register := new(bytes.Buffer)

		exitStatus := 0
		env := append(t.Env, registers...)

		switch t.Runner {
		case "sh":
			runner, err = run.NewShRunner(t.Code, t.Workdir, env)
		case "bash":
			runner, err = run.NewBashRunner(t.Code, t.Workdir, env)
		default:
			if t.Runner == "" {
				return errors.New("Runner not specified")
			}
			return fmt.Errorf("%s is not a supported runner", t.Runner)
		}
		if err != nil {
			return err
		}

		runOut, runErr, err := run.GetOutputs(runner)
		if err != nil {
			return err
		}

		if t.Register == "" {
			go io.Copy(stdout, runOut)
		} else {
			go io.Copy(register, runOut)
		}
		go io.Copy(stderr, runErr)

		startTime = time.Now()
		err = runner.Start()
		if err != nil {
			return err
		}

		err = runner.Wait()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); !ok {
					exitStatus = 1 // For unknown error exit codes set it to 1
				} else {
					exitStatus = waitStatus.ExitStatus()
				}
			}
			if ProvideExtraInfo {
				fmt.Printf("-- %s (%s) failed with exit status %d\n",
					t.Name, time.Since(startTime).String(), exitStatus)
			}
			return err
		}

		if t.Register != "" {
			r := fmt.Sprintf("%s=%s", t.Register, register.String())
			registers = append(registers, strings.TrimSpace(r))
		}

		if ProvideExtraInfo {
			fmt.Printf("-- %s (%s) finished with exit status %d\n",
				t.Name, time.Since(startTime).String(), exitStatus)
		}

	}
	return nil
}
