// Package errorchecks implements functions for inspecting error meanings.
//
// It is an implementation of Dave Cheney's proposition:
//     "Assert errors for behaviour, not type"
//
// To implement this, we assert for things that will guide the way we handle
// the error instead of asserting for the exact type of the error.
package errorchecks

import (
	"github.com/pkg/errors"
)

func cause(err error) error {
	return errors.Cause(err)
}
