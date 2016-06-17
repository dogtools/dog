package execute

import (
	"errors"
	"fmt"
	"os"

	"github.com/xsb/dog/types"
)

type runner struct {
	taskHierarchy map[string][]*types.Task
	eventsChan    chan *types.Event
}

func buildHierarchy(tm types.TaskMap) (map[string][]*types.Task, error) {
	th := make(map[string][]*types.Task, len(tm))

	for n, t := range tm {
		i := 0

		pres, ok := t.Pre.([]string)
		if !ok {
			return th, errors.New("Invalid pre directive")
		}

		posts, ok := t.Post.([]string)
		if !ok {
			return th, errors.New("Invalid post directive")
		}

		th[n] = make([]*types.Task, len(pres)+len(posts)+1)

		for _, preName := range pres {
			pre, found := tm[preName]
			if !found {
				return th, errors.New(
					"Task " + preName + " does not exist",
				)
			}
			th[n][i] = pre
			i++
		}

		th[n][i] = t
		i++

		for _, postName := range posts {
			post, found := tm[postName]
			if !found {
				return th, errors.New(
					"Task " + postName + " does not exist",
				)
			}
			th[n][i] = post
			i++
		}
	}

	return th, nil
}

func findCycles(th map[string][]*types.Task) error {
	for name, tasks := range th {
		for _, task := range tasks {
			if task.Name == name {
				continue
			}

			for _, subTask := range th[task.Name] {
				if subTask.Name == name {
					return errors.New("Hooks cycle in task " + name)
				}
			}
		}
	}
	return nil
}

// NewRunner creates a new runner that contains a list of all execution paths.
func NewRunner(tm types.TaskMap) (*runner, error) {
	th, err := buildHierarchy(tm)
	if err != nil {
		return nil, err
	}

	if err = findCycles(th); err != nil {
		return nil, err
	}

	return &runner{
		taskHierarchy: th,
		eventsChan:    make(chan *types.Event, 2048),
	}, nil
}

// Run executes the execution path for a given task.
func (r *runner) Run(taskName string) {
	tasks, found := r.taskHierarchy[taskName]
	if !found {
		fmt.Println("Task " + taskName + " does not exist")
		os.Exit(1)
	}

	executors := map[string]*Executor{}
	go func() {
		for _, t := range tasks {
			var e *Executor
			if t.Executor == "" {
				e = SystemExecutor
			} else {
				e, found = executors[t.Executor]
				if !found {
					e = NewExecutor(t.Executor)
					executors[t.Executor] = e
				}
			}

			if err := e.Exec(t, r.eventsChan); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()
	r.waitFor((tasks[len(tasks)-1]).Name)
}

func (r *runner) waitFor(taskName string) {
	for {
		select {
		case event := <-r.eventsChan:
			switch event.Name {
			case "start":
				fmt.Println(" - " + event.Task + " started")
			case "output":
				if body, ok := event.Extras["body"].([]byte); ok {
					fmt.Println(string(body))
				}
			case "end":
				if statusCode, ok := event.Extras["statusCode"].(int); ok {
					fmt.Println(
						fmt.Sprintf(" - %s finished with status code %d", event.Task, statusCode),
					)
					if statusCode != 0 || event.Task == taskName {
						os.Exit(statusCode)
					}
				} else {
					os.Exit(1)
				}
			}
		}
	}
}
