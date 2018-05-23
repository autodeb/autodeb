package oauth

import (
	"errors"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
)

const userIDSessionKey = "userid"

func (backend *backend) session(r *http.Request) *sessions.Session {
	session, _ := backend.sessionsManager.Get(r)
	return session
}

func (backend *backend) clearSession(r *http.Request, w http.ResponseWriter) {
	session := backend.session(r)
	session.Expire()
	session.Save(r, w)
}

func (backend *backend) getUserID(r *http.Request) (uint, error) {
	session := backend.session(r)

	userID, _ := session.Get(userIDSessionKey)
	if userID, ok := userID.(uint); ok {
		return userID, nil
	}

	return 0, errors.New("no userid in session")
}

func (backend *backend) setUserID(r *http.Request, w http.ResponseWriter, id uint) {
	session := backend.session(r)
	session.Set(userIDSessionKey, id)
	session.Save(r, w)
}
