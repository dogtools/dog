package execute

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dogtools/dog/types"
)

type runner struct {
	taskHierarchy map[string][]*types.Task
	eventsChan    chan *types.Event
	printFooter   bool
}

func isCyclic(chain []*types.Task) bool {
	maxLen := len(chain) / 2
	for i := 2; i < maxLen; i++ {
		a := chain[:i]
		b := chain[i : 2*i]
		for x, c := range a {
			if c != b[x] {
				return false
			}
		}
		return true
	}
	return false
}

func generateChainFor(t *types.Task, tm types.TaskMap, chain []*types.Task) ([]*types.Task, error) {
	var err error
	if isCyclic(chain) {
		return nil, errors.New("Task " + t.Name + " has a hook cycle")
	}

	for _, preName := range t.Pre {
		pre, found := tm[preName]
		if !found {
			return nil, errors.New(
				"Task " + preName + " does not exist",
			)
		}

		for _, prePre := range pre.Pre {
			if prePre == t.Name {
				return nil, errors.New("Task " + preName + " has a hook cycle")
			}
		}

		chain, err = generateChainFor(pre, tm, chain)
		if err != nil {
			return nil, err
		}
	}

	chain = append(chain, t)

	for _, postName := range t.Post {
		post, found := tm[postName]
		if !found {
			return nil, errors.New(
				"Task " + postName + " does not exist",
			)
		}
		chain, err = generateChainFor(post, tm, chain)
		if err != nil {
			return nil, err
		}
	}

	return chain, nil
}

func buildHierarchy(tm types.TaskMap) (map[string][]*types.Task, error) {
	th := make(map[string][]*types.Task, len(tm))

	for n, t := range tm {
		chain, err := generateChainFor(t, tm, []*types.Task{})
		if err != nil {
			return nil, err
		}
		th[n] = chain
	}

	return th, nil
}

func formatDuration(d time.Duration) (s string) {
	timeMsg := ""

	if d.Hours() > 1.0 {
		timeMsg += fmt.Sprintf("%1.0fh", d.Hours())
	}

	if d.Minutes() > 1.0 {
		timeMsg += fmt.Sprintf("%1.0fm", d.Minutes())
	}

	if d.Seconds() > 1.0 {
		timeMsg += fmt.Sprintf("%1.0fs", d.Seconds())
	} else {
		timeMsg += fmt.Sprintf("%1.3fs", d.Seconds())
	}

	return timeMsg
}

// NewRunner creates a new runner that contains a list of all execution paths.
func NewRunner(tm types.TaskMap, printFooter bool) (*runner, error) {
	th, err := buildHierarchy(tm)
	if err != nil {
		return nil, err
	}

	return &runner{
		taskHierarchy: th,
		eventsChan:    make(chan *types.Event, 2048),
		printFooter:   printFooter,
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
			case "output":
				if body, ok := event.Extras["body"].([]byte); ok {
					fmt.Println(string(body))
				}
			case "end":
				if statusCode, ok := event.Extras["statusCode"].(int); ok {
					if elapsed, ok := event.Extras["elapsed"].(time.Duration); ok {
						if r.printFooter {
							fmt.Printf("-- %s took %s and finished with status code %d\n", event.Task, formatDuration(elapsed), statusCode)
						}

						if statusCode != 0 || event.Task == taskName {
							os.Exit(statusCode)
						}
					}

				} else {
					os.Exit(1)
				}
			}
		}
	}
}
