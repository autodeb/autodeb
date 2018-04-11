package decorators

import (
	"net/http"
)

// HTMLHeaders adds the Content-Type text/html header to the request.
// It adds the header before running the handler. Otherwise, it would be too
// late to add headers.
func HTMLHeaders(h http.HandlerFunc) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h(w, r)
	}
	return handler
}
