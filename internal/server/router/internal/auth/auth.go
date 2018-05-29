package auth

import (
	"net/http"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// UserHandlerFunc is net/http HandlerFunc with an additional User parameter
type UserHandlerFunc = func(http.ResponseWriter, *http.Request, *models.User)

func identifyUser(appCtx *appctx.Context, r *http.Request) (*models.User, error) {
	// There is two possible methods of authentication.
	// 1. Using an access token
	// 2. Trough the auth backend

	// 1. Using an access token
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		if splitAuthHeader := strings.Split(authHeader, " "); len(splitAuthHeader) >= 2 {
			token := splitAuthHeader[1]

			user, err := appCtx.TokensService().GetUserByToken(token)
			if err != nil {
				return nil, err
			}
			if user != nil {
				return user, nil
			}
		}
	}

	// 2. Trough the auth backend
	if user, err := appCtx.AuthBackend().GetUser(r); err != nil {
		return nil, err
	} else if user != nil {
		return user, nil
	}

	return nil, nil
}

// MaybeWithUser retrieves the connected user and calls the provided function
// with the user. If the user is not connected, user is nil.
func MaybeWithUser(fn UserHandlerFunc, appCtx *appctx.Context) http.Handler {

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Attempt to identify the user
		user, err := identifyUser(appCtx, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
