package oauth

import (
	"net/http"
	"strconv"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
)

const authBackendUserIDSessionKey = "auth-backend-user-id"

func (backend *backend) session(r *http.Request) *sessions.Session {
	session, _ := backend.sessionsManager.Get(r)
	return session
}

func (backend *backend) clearSession(r *http.Request, w http.ResponseWriter) {
	session := backend.session(r)
	session.Expire()
	session.Save(r, w)
}

func (backend *backend) getAuthBackendUserID(r *http.Request) (uint, error) {
	session := backend.session(r)

	authBackendUserIDString, err := session.Get(authBackendUserIDSessionKey)
	if err != nil {
		return 0, err
	}

	authBackendUserID, err := strconv.Atoi(authBackendUserIDString)

	return uint(authBackendUserID), errors.WithMessage(err, "could not cast user id to int")
}

func (backend *backend) setUserID(r *http.Request, w http.ResponseWriter, id uint) {
	session := backend.session(r)
	session.Set(
		authBackendUserIDSessionKey,
		strconv.Itoa(int(id)),
	)
	session.Save(r, w)
}
