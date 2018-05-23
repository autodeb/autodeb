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
	cookie := w.Header()["Set-Cookie"][0]
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	request.Header.Add("Cookie", cookie)
	return request
}

func TestFlash(t *testing.T) {
	manager := setupTest()

	// Initial request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	response := httptest.NewRecorder()

	// Add flashes
	session, _ := manager.Get(request)
	session.Flash("error", "this is an error")
	session.Flash("error", "this is a second error")
	session.Flash("info", "this is information")
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
