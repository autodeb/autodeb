package oauth

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

func (service *service) oauthConfigWithRedirectURL() *oauth2.Config {
	cfg := service.oauthProvider.Config()

	url := strings.Trim(service.serverURL, "/")
	url = url + "/auth/callback"

	cfg.RedirectURL = url

	return cfg
}

func (service *service) LoginHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		oauthCfg := service.oauthConfigWithRedirectURL()

		url := oauthCfg.AuthCodeURL("")

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	}

	return http.HandlerFunc(handlerFunc)
}

func (service *service) LogoutHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		service.logout(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(handlerFunc)
}

func (service *service) AuthHandler() http.Handler {
	//TODO: setup a router so that we don't answer to all routes in /auth/
	return service.callbackHandler()
}

func (service *service) callbackHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		authCode := r.FormValue("code")

		oauthCfg := service.oauthConfigWithRedirectURL()

		token, err := oauthCfg.Exchange(oauth2.NoContext, authCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, username, err := service.oauthProvider.UserInfo(token.AccessToken)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := service.login(r, w, userID, username); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(handlerFunc)
}
