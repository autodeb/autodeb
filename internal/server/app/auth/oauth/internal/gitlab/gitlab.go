package gitlab

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

//OAuthProvider provides oauth config for gitlab
type OAuthProvider struct {
	baseURL      *url.URL
	clientID     string
	clientSecret string
}

//New returns a gitlab oauth provider
func New(baseURL *url.URL, clientID, clientSecret string) *OAuthProvider {
	oauthProvider := &OAuthProvider{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
	return oauthProvider
}

//Config returns the gitlab oauth2 config
func (provider *OAuthProvider) Config() *oauth2.Config {
	authURL := provider.baseURL.ResolveReference(
		&url.URL{Path: "/oauth/authorize"},
	).String()

	tokenURL := provider.baseURL.ResolveReference(
		&url.URL{Path: "/oauth/token"},
	).String()

	cfg := &oauth2.Config{
		ClientID:     provider.clientID,
		ClientSecret: provider.clientSecret,
		Scopes:       []string{"read_user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	return cfg
}

//UserInfo returns the user id and username
func (provider *OAuthProvider) UserInfo(authToken string) (uint, string, error) {
	userURL := provider.baseURL.ResolveReference(
		&url.URL{Path: "/api/v4/user"},
	).String()
	userURL = userURL + "?access_token=" + authToken

	response, err := http.Get(userURL)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, "", err
	}

	var apiUser struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	if err := json.Unmarshal(contents, &apiUser); err != nil {
		return 0, "", err
	}

	return apiUser.ID, apiUser.Username, nil
}
