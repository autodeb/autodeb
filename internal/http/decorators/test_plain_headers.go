package decorators

import (
	"net/http"
)

// TextPlainHeaders adds the Content-Type text/plain header to the request.
// It adds the header before running the handler. Otherwise, it would be too
// late to add headers.
func TextPlainHeaders(h http.HandlerFunc) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		h(w, r)
	}
	return handler
}
