package execute

import (
	"errors"
	"fmt"
	"os"

	"github.com/dogtools/dog/types"
)

type runner struct {
	taskHierarchy map[string][]*types.Task
	eventsChan    chan *types.Event
}

func generateChainFor(t *types.Task, tm types.TaskMap) ([]*types.Task, error) {
	chain := []*types.Task{}
	if pres, ok := t.Pre.([]string); ok {
		for _, preName := range pres {
			pre, found := tm[preName]
			if !found {
				return nil, errors.New(
					"Task " + preName + " does not exist",
				)
			}
			preChain, err := generateChainFor(pre, tm)
			if err != nil {
				return nil, err
			}
			chain = append(chain, preChain...)
		}
	}

	chain = append(chain, t)

	if posts, ok := t.Post.([]string); ok {
		for _, postName := range posts {
			post, found := tm[postName]
			if !found {
				return []*types.Task{}, errors.New(
					"Task " + postName + " does not exist",
				)
			}
			postChain, err := generateChainFor(post, tm)
			if err != nil {
				return nil, err
			}
			chain = append(chain, postChain...)
		}
	}

	return chain, nil
}

func buildHierarchy(tm types.TaskMap) (map[string][]*types.Task, error) {
	th := make(map[string][]*types.Task, len(tm))

	for n, t := range tm {
		chain, err := generateChainFor(t, tm)
		if err != nil {
			return nil, err
		}
		th[n] = chain
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
