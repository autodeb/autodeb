//Package sessions is wrapper around github.com/gorilla/sessions.
//
// Dfferences include:
//  - keys and values can only be strings
//  - way less exposed API surface
package sessions

import (
	"net/http"

	gorillaSessions "github.com/gorilla/sessions"
)

//Manager is a sessions manager. It wraps a session store to ensure that all
//callers use the same session name.
type Manager struct {
	store gorillaSessions.Store
	name  string
}

//NewManager creates a new sessions manager
func NewManager(store gorillaSessions.Store, name string) *Manager {
	manager := &Manager{
		store: store,
		name:  name,
	}
	return manager
}

//Get returns the session associated with the given request
func (m *Manager) Get(r *http.Request) (*Session, error) {
	gorillaSession, err := m.store.Get(r, m.name)
	session := newSession(gorillaSession)
	return session, err
}
