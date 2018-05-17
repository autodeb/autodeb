package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const userIDSessionKey = "userid"

func (service *Service) session(r *http.Request) *sessions.Session {
	session, _ := service.sessionStore.Get(r, "autodeb")
	return session
}

func (service *Service) clearSession(r *http.Request, w http.ResponseWriter) {
	session := service.session(r)
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func (service *Service) getUserID(r *http.Request) (uint, error) {
	session := service.session(r)

	if userID, ok := session.Values[userIDSessionKey].(uint); ok {
		return userID, nil
	}

	return 0, errors.New("no userid in session")
}

func (service *Service) setUserID(r *http.Request, w http.ResponseWriter, id uint) {
	session := service.session(r)
	session.Values[userIDSessionKey] = id
	session.Save(r, w)
}
