//Package disabled implements a disabled authentification backend.
//
//Authentification pages display a message indicating that authentification
//been disabled.
package disabled

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

type backend struct{}

//NewBackend returns a new disabled authentificiation backend
func NewBackend() auth.Backend {
	backend := &backend{}
	return backend
}

func handler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "The administrators of this server have not enabled authentificaiton.")
	}
	return http.HandlerFunc(handlerFunc)
}

func (backend *backend) LoginHandler() http.Handler {
	return handler()
}

func (backend *backend) LogoutHandler() http.Handler {
	return handler()
}

func (backend *backend) AuthHandler() http.Handler {
	return handler()
}

func (backend *backend) GetUser(r *http.Request) (*models.User, error) {
	return nil, nil
}
