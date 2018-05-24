//Package oauth implements an authentification backend based on OAuth. It
//allows for authenticating against an OAuth2 provider.
package oauth

import (
	"net/http"
	"net/url"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

type backend struct {
	db              *database.Database
	sessionsManager *sessions.Manager
	oauthProvider   Provider
	serverURL       *url.URL
}

// NewBackend returns a new oauth authentification backend
func NewBackend(db *database.Database, sessionsManager *sessions.Manager, oauthProvider Provider, serverURL *url.URL) auth.Backend {
	backend := &backend{
		db:              db,
		sessionsManager: sessionsManager,
		oauthProvider:   oauthProvider,
		serverURL:       serverURL,
	}
	return backend
}

func (backend *backend) GetUser(r *http.Request) (*models.User, error) {
	id, err := backend.getUserID(r)
	if err != nil {
		// Ignore the error, there is no user in this session
		return nil, nil
	}

	// Get the user model
	user, err := backend.db.GetUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (backend *backend) SetUser(r *http.Request, w http.ResponseWriter, user *models.User) {
	backend.setUserID(r, w, user.ID)
}

func (backend *backend) login(r *http.Request, w http.ResponseWriter, id uint, username string) error {
	user, err := backend.db.GetUser(id)
	if err != nil {
		return err
	}

	if user == nil {
		user, err = backend.db.CreateUser(id, username)
		if err != nil {
			return err
		}
	}

	backend.SetUser(r, w, user)

	return nil
}

func (backend *backend) logout(r *http.Request, w http.ResponseWriter) {
	backend.clearSession(r, w)
}
