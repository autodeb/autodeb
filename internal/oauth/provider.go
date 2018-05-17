package oauth

import (
	"fmt"
	"net/url"

	"golang.org/x/oauth2"

	"salsa.debian.org/autodeb-team/autodeb/internal/oauth/internal/gitlab"
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
	switch provider {
	case "gitlab":
		provider := gitlab.New(baseURL, clientID, clientSecret)
		return provider, nil
	default:
		return nil, fmt.Errorf("unknown provider %s", provider)
	}
}
