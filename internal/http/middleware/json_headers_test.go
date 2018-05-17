package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"

	"github.com/stretchr/testify/assert"
)

func TestJSONHeaders(t *testing.T) {
	// Create an empty handler
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"hello": "world"}`)
	}
	handler := http.Handler(http.HandlerFunc(handlerFunc))

	// Test the empty handler
	request, _ := http.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType := responseRecorder.Result().Header.Get("Content-Type")
	assert.NotEqual(t, "application/json", contentType)

	// Weap the handler
	handler = middleware.JSONHeaders(handler)

	// Test the decorated handler
	request, _ = http.NewRequest("GET", "/", nil)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType = responseRecorder.Result().Header.Get("Content-Type")
	assert.Equal(t, "application/json", contentType)
}
