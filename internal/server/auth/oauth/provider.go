package oauth

import (
	"net/url"

	"golang.org/x/oauth2"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth/oauth/internal/gitlab"
)

// Provider represents an oauth provider. It provides a config and a method to
// obtain user information from an access token.
type Provider interface {
	Config() *oauth2.Config

	//UserInfo returns the user ID and username given an auth token
	UserInfo(token string) (uint, string, error)
}

//NewProvider creates a new OAuth Provider
func NewProvider(provider string, baseURL *url.URL, clientID, clientSecret string) (Provider, error) {
	if baseURL.String() == "" {
		return nil, errors.New("baseURL cannot be empty")
	}
	if clientID == "" {
		return nil, errors.New("clientID cannot be empty")
	}
	if clientSecret == "" {
		return nil, errors.New("clientSecret cannot be empty")
	}

	switch provider {
	case "gitlab":
		provider := gitlab.New(baseURL, clientID, clientSecret)
		return provider, nil
	default:
		return nil, errors.Errorf("unknown provider %s", provider)
	}
}
