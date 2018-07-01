package exec

import (
	osexec "os/exec"
	"syscall"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

// ExitCodeFromError attempts to find the exit code from an exec error
func ExitCodeFromError(err error) (int, error) {

	// Retrieve the error's underlying cause
	cause := errors.Cause(err)

	exitError, ok := cause.(*osexec.ExitError)
	if !ok {
		return -1, errors.Errorf("could not cast %+v into ExitError", cause)
	}

	waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
	if !ok {
		return -1, errors.Errorf("could not cast %+v into WaitStatus", exitError.Sys())
	}

	return waitStatus.ExitStatus(), nil
}
