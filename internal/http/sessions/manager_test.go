package sessions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gorillaSessions "github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
)

func setupTest() *sessions.Manager {
	store := gorillaSessions.NewCookieStore(
		[]byte("test-cookie-secret"),
	)
	manager := sessions.NewManager(store, "test-session-name")
	return manager
}

func getRequestWithResponseCookies(w http.ResponseWriter) *http.Request {
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	cookies := w.Header()["Set-Cookie"]
	request.Header.Add("Cookie", cookies[len(cookies)-1])
	return request
}

func TestFlash(t *testing.T) {
	manager := setupTest()

	// Initial request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	response := httptest.NewRecorder()

	// Add flashes trough the session
	session, _ := manager.Get(request)
	session.Flash("error", "this is an error")
	session.Flash("error", "this is a second error")
	session.Save(request, response)

	// Add flashes trough the manager
	manager.Flash(request, response, "info", "this is information")
	session.Save(request, response)

	// Second request
	request = getRequestWithResponseCookies(response)
	session, _ = manager.Get(request)

	// Get flashes
	flashes := session.Flashes()

	assert.Equal(t, 2, len(flashes))
	assert.Equal(t, []string{"this is an error", "this is a second error"}, flashes["error"])
	assert.Equal(t, []string{"this is information"}, flashes["info"])
}
