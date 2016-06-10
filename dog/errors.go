package dog

import "errors"

// AlreadyRegisteredErr indicates that the executor has been registered twice.
var AlreadyRegisteredErr = errors.New("You are trying to register an executor twice")
