package dog

// DefaultRunner defines the runner to use in case the task does not specify it.
var DefaultRunner = "sh"

// ProvideExtraInfo specifies if dog needs to provide execution info (duration,
// exit status) after task execution.
var ProvideExtraInfo bool

// deprecation warning flags
var deprecationWarningRun bool
var deprecationWarningExec bool
