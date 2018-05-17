package auth

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

func oauthConfigWithRedirectURL(app *app.App) *oauth2.Config {
	cfg := app.OAuthProvider().Config()

	url := strings.Trim(app.Config().ServerURL, "/")
	url = url + "/auth/callback"

	cfg.RedirectURL = url

	return cfg
}

//LoginGetHandler returns a handler that redirects to the oauth provider
func LoginGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		oauthCfg := oauthConfigWithRedirectURL(app)

		url := oauthCfg.AuthCodeURL("")

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	}

	return http.HandlerFunc(handler)
}

//CallbackGetHandler returns a handlers that handles the oauth callback
func CallbackGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		oauthCfg := oauthConfigWithRedirectURL(app)

		authCode := r.FormValue("code")
		token, err := oauthCfg.Exchange(oauth2.NoContext, authCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, _, err := app.OAuthProvider().UserInfo(token.AccessToken)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: create user in our database

		session := app.Session(r)
		setUserID(session, userID)
		session.Save(r, w)

		http.Redirect(w, r, "/account", http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(handler)
}
