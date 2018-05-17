package auth

import (
	"errors"

	"github.com/gorilla/sessions"
)

const userIDSessionKey = "userid"

func getUserID(s *sessions.Session) (int, error) {
	if userID, ok := s.Values[userIDSessionKey].(int); ok {
		return userID, nil
	}
	return 0, errors.New("no userid in session")
}

func setUserID(s *sessions.Session, id int) {
	s.Values[userIDSessionKey] = id
}
