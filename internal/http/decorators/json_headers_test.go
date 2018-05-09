package decorators_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"

	"github.com/stretchr/testify/assert"
)

func TestJSONHeaders(t *testing.T) {
	// Create an empty handler
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"hello": "world"}`)
	}
	handler := http.HandlerFunc(handlerFunc)

	// Test the empty handler
	request, _ := http.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType := responseRecorder.Result().Header.Get("Content-Type")
	assert.NotEqual(t, "application/json", contentType)

	// Decorate the handler
	handlerFunc = decorators.JSONHeaders(handlerFunc)
	handler = http.HandlerFunc(handlerFunc)

	// Test the decorated handler
	request, _ = http.NewRequest("GET", "/", nil)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType = responseRecorder.Result().Header.Get("Content-Type")
	assert.Equal(t, "application/json", contentType)
}
