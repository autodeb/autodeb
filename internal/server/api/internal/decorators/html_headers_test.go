package decorators_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/api/internal/decorators"
)

func TestHTMLHeaders(t *testing.T) {
	// Create an empty handler
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `Hello World`)
	}
	handler := http.HandlerFunc(handlerFunc)

	// Test the empty handler
	request, _ := http.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType := responseRecorder.Result().Header.Get("Content-Type")
	if contentType == "application/json" {
		t.Error("Content-Type should not be application/json yet")
	}

	// Decoreate the handler
	handlerFunc = decorators.HTMLHeaders(handlerFunc)
	handler = http.HandlerFunc(handlerFunc)

	// Test the decorated handler
	request, _ = http.NewRequest("GET", "/", nil)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	contentType = responseRecorder.Result().Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Error("Content-Type should be application/json: ", contentType)
	}
}
