package oauth

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

func (backend *backend) oauthConfigWithRedirectURL() *oauth2.Config {
	cfg := backend.oauthProvider.Config()

	url := strings.Trim(backend.serverURL.String(), "/")
	url = url + "/auth/callback"

	cfg.RedirectURL = url

	return cfg
}

func (backend *backend) LoginHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		oauthCfg := backend.oauthConfigWithRedirectURL()

		url := oauthCfg.AuthCodeURL("")

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	}

	return http.HandlerFunc(handlerFunc)
}

func (backend *backend) LogoutHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		backend.logout(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(handlerFunc)
}

func (backend *backend) AuthHandler() http.Handler {
	//TODO: setup a router so that we don't answer to all routes in /auth/
	return backend.callbackHandler()
}

func (backend *backend) callbackHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		authCode := r.FormValue("code")

		oauthCfg := backend.oauthConfigWithRedirectURL()

		token, err := oauthCfg.Exchange(oauth2.NoContext, authCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, username, err := backend.oauthProvider.UserInfo(token.AccessToken)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := backend.login(r, w, userID, username); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(handlerFunc)
}
