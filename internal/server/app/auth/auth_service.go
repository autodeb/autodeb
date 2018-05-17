package auth

import (
	"net/http"

	"github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// Service manages authentification.
type Service struct {
	db           *database.Database
	sessionStore sessions.Store
}

// NewService returns a new suthentification service
func NewService(db *database.Database, sessionStore sessions.Store) *Service {
	service := &Service{
		db:           db,
		sessionStore: sessionStore,
	}

	return service
}

// GetUser retrieves the user associated with a request
func (service *Service) GetUser(r *http.Request) (*models.User, error) {
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

// SetUser associates a user with a request
func (service *Service) SetUser(r *http.Request, w http.ResponseWriter, user *models.User) {
	service.setUserID(r, w, user.ID)
}

// Login finds or creates a user and associates it with the request
func (service *Service) Login(r *http.Request, w http.ResponseWriter, id uint, username string) error {
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

// Logout logs out the user
func (service *Service) Logout(r *http.Request, w http.ResponseWriter) {
	service.clearSession(r, w)
}
