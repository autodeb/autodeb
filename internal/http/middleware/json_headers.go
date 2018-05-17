package middleware

import (
	"net/http"
)

// JSONHeaders adds the Content-Type application/json header to the request.
// It adds the header before running the handler. Otherwise, it would be too
// late to add headers.
func JSONHeaders(h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}
