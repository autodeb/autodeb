package oauth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const userIDSessionKey = "userid"

func (backend *backend) session(r *http.Request) *sessions.Session {
	session, _ := backend.sessionStore.Get(r, "autodeb")
	return session
}

func (backend *backend) clearSession(r *http.Request, w http.ResponseWriter) {
	session := backend.session(r)
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func (backend *backend) getUserID(r *http.Request) (uint, error) {
	session := backend.session(r)

	if userID, ok := session.Values[userIDSessionKey].(uint); ok {
		return userID, nil
	}

	return 0, errors.New("no userid in session")
}

func (backend *backend) setUserID(r *http.Request, w http.ResponseWriter, id uint) {
	session := backend.session(r)
	session.Values[userIDSessionKey] = id
	session.Save(r, w)
}
