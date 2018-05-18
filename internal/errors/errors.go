//Package errors exposes a subset of the github.com/pkg/errors API.
package errors

import (
	"fmt"
	pkgerrors "github.com/pkg/errors"
)

//New creates a new error
func New(message string) error {
	return pkgerrors.New(message)
}

//Errorf creates a new formatted error
func Errorf(format string, args ...interface{}) error {
	return pkgerrors.Errorf(format, args...)
}

//WithMessage annotates an error with a message. Unlike github.com/pkg/errors,
//this implementation also annotates the error with a stack trace if there
//is not already one.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	type stackTracer interface {
		StackTrace() pkgerrors.StackTrace
	}

	// Try to find a StackTracer down the error cause chain.
	current := err
	for current != nil {
		if _, ok := current.(stackTracer); ok {
			break
		}
		current = nextCause(current)
	}

	// The error is not yet annotated with a stack trace, Wrap it.
	if current == nil {
		err = pkgerrors.Wrap(err, message)
	}

	return pkgerrors.WithMessage(err, message)
}

//WithMessagef annonates an error with a formatted message
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return WithMessage(err, message)
}

//nextCause returns the direct cause of the error
func nextCause(err error) error {
	type causer interface {
		Cause() error
	}

	if causer, ok := err.(causer); ok {
		return causer.Cause()
	}

	return nil
}
