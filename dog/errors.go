package dog

import "errors"

// ErrAlreadyRegistered indicates that the executor has been registered twice.
var ErrAlreadyRegistered = errors.New("You are trying to register an executor twice")
