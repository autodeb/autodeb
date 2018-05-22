package appctxtest

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth/disabled"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// fakeAuthBackend allows to preset return values of GetUser to facilitate
// testing.
type fakeAuthBackend struct {
	auth.Backend

	//Return values of GetUser
	User  *models.User
	Error error
}

func newFakeAuthBackend() *fakeAuthBackend {
	backend := &fakeAuthBackend{
		Backend: disabled.NewBackend(),
		User:    nil,
		Error:   nil,
	}
	return backend
}

//GetUser returns the preset user instead of doing any authentication work.
func (backend *fakeAuthBackend) GetUser(r *http.Request) (*models.User, error) {
	return backend.User, backend.Error
}
