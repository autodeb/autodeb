package config

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/net/url"
)

// ParseConfig parses a configuration file to create a server config
func ParseConfig(filepath string, fs filesystem.FS) (*Config, error) {
	file, err := fs.Open(filepath)
	if err != nil {
		return nil, errors.WithMessage(err, "could not open configuration file")
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.WithMessage(err, "could not read config file")
	}

	// Default ServerURL
	serverURL := &url.URL{}
	if err := serverURL.UnmarshalBinary([]byte("http://localhost:8071")); err != nil {
		return nil, errors.WithMessage(err, "cannot parse default server url")
	}

	// Default repositories url
	aptlyRepositoryBaseURL := &url.URL{}
	if err := aptlyRepositoryBaseURL.UnmarshalBinary([]byte("http://localhost:8071/repos")); err != nil {
		return nil, errors.WithMessage(err, "cannot parse default repositories url")
	}

	// Default aptly api url
	aptlyAPIURL := &url.URL{}
	if err := aptlyAPIURL.UnmarshalBinary([]byte("http://localhost:8080/api")); err != nil {
		return nil, errors.WithMessage(err, "cannot parse default repositories url")
	}

	// Create the config, with defaults
	config := &Config{
		DB: &DBConfig{
			Driver:           "sqlite3",
			ConnectionString: "database.sqlite",
		},
		HTTP: &HTTPServerConfig{
			Address: ":8071",
		},
		Auth: &AuthConfig{
			AuthentificationBackend: "disabled",
			OAuth: &OAuthConfig{
				Provider:     "gitlab",
				BaseURL:      "https://salsa.debian.org",
				ClientID:     "",
				ClientSecret: "",
			},
		},
		Aptly: &Aptly{
			APIURL:            aptlyAPIURL,
			RepositoryBaseURL: aptlyRepositoryBaseURL,
		},
		ServerURL:             serverURL,
		DataDirectory:         "data",
		TemplatesDirectory:    "web/templates",
		StaticFilesDirectory:  "web/static",
		TemplatesCacheEnabled: true,
		LogLevel:              log.InfoLevel,
	}

	if metadata, err := toml.DecodeReader(bytes.NewReader(b), &config); err != nil {
		return nil, errors.WithMessage(err, "could not decode configuration file")
	} else if keys := metadata.Undecoded(); len(keys) > 0 {
		return nil, errors.Errorf("unrecognized configuration key: %s", keys[0].String())
	}

	return config, nil
}
