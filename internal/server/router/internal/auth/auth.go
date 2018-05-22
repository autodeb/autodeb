package auth

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// UserHandlerFunc is net/http HandlerFunc with an additional User parameter
type UserHandlerFunc = func(http.ResponseWriter, *http.Request, *models.User)

// MaybeWithUser retrieves the connected user and calls the provided function
// with the user. If the user is not connected, user is nil.
func MaybeWithUser(fn UserHandlerFunc, appCtx *appctx.Context) http.Handler {

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get the user
		user, err := appCtx.AuthBackend().GetUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Call the function
		fn(w, r, user)

	}

	return http.HandlerFunc(handlerFunc)
}

// WithUser promises to call the provided function with a non-nil user.
// Other requests are redirected to the authentification page.
func WithUser(fn UserHandlerFunc, appCtx *appctx.Context) http.Handler {

	userHandlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		// Redirect if the user is not authenticated
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Call the function
		fn(w, r, user)
	}

	handler := MaybeWithUser(userHandlerFunc, appCtx)

	return handler
}
