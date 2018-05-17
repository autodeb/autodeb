package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
)

func TestHTMLHeaders(t *testing.T) {
	// Create an empty handler
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `Hello World`)
	}
	handler := http.Handler(http.HandlerFunc(handlerFunc))

	// Test the empty handler
	request, _ := http.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType := responseRecorder.Result().Header.Get("Content-Type")
	if contentType == "application/json" {
		t.Error("Content-Type should not be application/json yet")
	}

	// Wrap the handler
	handler = middleware.HTMLHeaders(handler)

	// Test the decorated handler
	request, _ = http.NewRequest("GET", "/", nil)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType = responseRecorder.Result().Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Error("Content-Type should be application/json: ", contentType)
	}
}
