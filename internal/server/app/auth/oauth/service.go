package oauth

import (
	"net/http"

	"github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

type service struct {
	db            *database.Database
	sessionStore  sessions.Store
	oauthProvider Provider
	serverURL     string
}

// NewService returns a new oauth authentification service
func NewService(db *database.Database, sessionStore sessions.Store, oauthProvider Provider, serverURL string) auth.Service {
	service := &service{
		db:            db,
		sessionStore:  sessionStore,
		oauthProvider: oauthProvider,
		serverURL:     serverURL,
	}
	return service
}

func (service *service) GetUser(r *http.Request) (*models.User, error) {
	id, err := service.getUserID(r)
	if err != nil {
		// Ignore the error, there is no user in this session
		return nil, nil
	}

	// Get the user model
	user, err := service.db.GetUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *service) SetUser(r *http.Request, w http.ResponseWriter, user *models.User) {
	service.setUserID(r, w, user.ID)
}

func (service *service) login(r *http.Request, w http.ResponseWriter, id uint, username string) error {
	user, err := service.db.GetUser(id)
	if err != nil {
		return err
	}

	if user == nil {
		user, err = service.db.CreateUser(id, username)
		if err != nil {
			return err
		}
	}

	service.SetUser(r, w, user)

	return nil
}

func (service *service) logout(r *http.Request, w http.ResponseWriter) {
	service.clearSession(r, w)
}
