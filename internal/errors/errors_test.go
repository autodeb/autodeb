package errors_test

import (
	goerrors "errors"
	"fmt"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"

	"github.com/stretchr/testify/assert"
)

func TestWithMessage(t *testing.T) {
	// stdlib errors should have no stack trace
	err := goerrors.New("test")
	assert.Equal(t, 0, strings.Count(fmt.Sprintf("%+v", err), "testing.go"))

	// Add a message to the error and it should contain a stack trace
	err = errors.WithMessage(err, "added message")
	assert.Equal(t, 1, strings.Count(fmt.Sprintf("%+v", err), "testing.go"))

	// Add another message to it and it should contain only one stack trace
	err = errors.WithMessage(err, "added message")
	assert.Equal(t, 1, strings.Count(fmt.Sprintf("%+v", err), "testing.go"))
}
