package oauth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const userIDSessionKey = "userid"

func (service *service) session(r *http.Request) *sessions.Session {
	session, _ := service.sessionStore.Get(r, "autodeb")
	return session
}

func (service *service) clearSession(r *http.Request, w http.ResponseWriter) {
	session := service.session(r)
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func (service *service) getUserID(r *http.Request) (uint, error) {
	session := service.session(r)

	if userID, ok := session.Values[userIDSessionKey].(uint); ok {
		return userID, nil
	}

	return 0, errors.New("no userid in session")
}

func (service *service) setUserID(r *http.Request, w http.ResponseWriter, id uint) {
	session := service.session(r)
	session.Values[userIDSessionKey] = id
	session.Save(r, w)
}
