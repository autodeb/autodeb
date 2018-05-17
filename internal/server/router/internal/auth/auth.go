package auth

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

//LoginGetHandler returns a handler that redirects to the oauth provider
func LoginGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		oauthCfg := app.OAuthProvider().Config()
		oauthCfg.RedirectURL = "http://localhost:8071/auth/callback"

		// Get the url to redirect to
		url := oauthCfg.AuthCodeURL("")

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(handler)
}

//CallbackGetHandler returns a handlers that handles the oauth callback
func CallbackGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		// Get the oauth config and set the redirect url
		oauthCfg := app.OAuthProvider().Config()
		oauthCfg.RedirectURL = "http://localhost:8071/auth/callback"

		// Obtain the auth token from the oauth provider
		authCode := r.FormValue("code")
		token, err := oauthCfg.Exchange(oauth2.NoContext, authCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Obtain the user information from the oauth provider
		userID, username, err := app.OAuthProvider().UserInfo(token.AccessToken)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "hello")
		fmt.Fprintf(w, "id: %d\n", userID)
		fmt.Fprintf(w, "username: %s\n", username)
	}

	return http.HandlerFunc(handler)
}
