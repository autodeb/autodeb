package sessions

import (
	"net/http"
	"strings"

	gorillaSessions "github.com/gorilla/sessions"
)

//Session stores values
type Session struct {
	gorillaSession *gorillaSessions.Session
}

func newSession(gorillaSession *gorillaSessions.Session) *Session {
	session := &Session{
		gorillaSession: gorillaSession,
	}
	return session
}

// Get will return the value associated with a key
func (s *Session) Get(key string) (interface{}, bool) {
	value, ok := s.gorillaSession.Values[key]
	return value, ok
}

// Set will set the value of a key
func (s *Session) Set(key string, value interface{}) {
	s.gorillaSession.Values[key] = value
}

// Keys returns available session keys
func (s *Session) Keys() []string {
	var keys []string
	for key := range s.gorillaSession.Values {
		if key, ok := key.(string); ok {
			keys = append(keys, key)
		}
	}
	return keys
}

// Expire will expire the session
func (s *Session) Expire() {
	s.gorillaSession.Options.MaxAge = -1
}

// Save will save the session
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return s.gorillaSession.Save(r, w)
}

const flashPrefix = "_flash_"

//Flash will add a message to a category of flashes
func (s *Session) Flash(r *http.Request, w http.ResponseWriter, category, message string) {
	s.gorillaSession.AddFlash(message, flashPrefix+category)
}

//Flashes will return all flashes and removes them from the request
func (s *Session) Flashes(r *http.Request) map[string][]string {
	flashes := make(map[string][]string)

	// Browse all session keys
	for _, key := range s.Keys() {
		// Verify that the key indicates that it is a flash
		if !strings.HasPrefix(key, flashPrefix) {
			continue
		}

		// Get the flash category
		category := strings.TrimPrefix(key, flashPrefix)

		// Get the flashes of this category, deleting them from the session
		values := s.gorillaSession.Flashes(flashPrefix + category)

		// Transform the values in message strings
		var messages []string
		for _, message := range values {
			if message, ok := message.(string); ok {
				messages = append(messages, message)
			}
		}

		if len(messages) > 0 {
			flashes[category] = messages
		}

	}

	return flashes
}
