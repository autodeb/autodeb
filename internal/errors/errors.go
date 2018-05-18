//Package errors exposes a subset of the github.com/pkg/errors API.
package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

//New creates a new error
func New(message string) error {
	return errors.New(message)
}

//Errorf creates a new formatted error
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

//WithMessage annotates an error with a message
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

//WithMessagef annonates an error with a formatted message
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return WithMessage(err, message)
}
