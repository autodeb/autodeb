//Package auth contains implementations of authentification backends
package auth

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//Backend provides everything that is needed to perform authentification
type Backend interface {
	LoginHandler() http.Handler
	LogoutHandler() http.Handler
	AuthHandler() http.Handler
	GetUser(*http.Request) (*models.User, error)
}
