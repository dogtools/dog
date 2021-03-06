package dog

// Task represents a task described in the Dogfile format.
type Task struct {
	// Name of the task.
	Name string

	// Description of the task.
	Description string

	// The code that will be executed.
	Code string

	// Defaults to operating system main shell.
	Runner string

	// Pre-hooks execute other tasks before starting the current one.
	Pre []string

	// Post-hooks are analog to pre-hooks but they are executed after
	// current task finishes its execution.
	Post []string

	// Default values for environment variables can be provided in the Dogfile.
	// They can be modified at execution time.
	Env []string

	// Sets the working directory for the task. Relative paths are
	// considered relative to the location of the Dogfile.
	Workdir string

	// Register stores the output of the task so it can be accessed by
	// other tasks in the task chain.
	//
	// When present, a new environment variable is injected in future
	// task chain runners using the register name as key and the output
	// as value.
	Register string
}
