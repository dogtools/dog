package dog

import "runtime"

// DefaultRunner defines the runner to use in case the task does not specify it.
//
// The value is automatically assigned based on the operating system when the
// package initializes.
var DefaultRunner string

// ProvideExtraInfo specifies if dog needs to provide execution info (duration,
// exit status) after task execution.
var ProvideExtraInfo bool

// deprecation warning flags
var deprecationWarningRun bool
var deprecationWarningExec bool

func init() {
	if runtime.GOOS == "windows" {
		DefaultRunner = "cmd" // not implemented yet
	} else {
		DefaultRunner = "sh"
	}
}
