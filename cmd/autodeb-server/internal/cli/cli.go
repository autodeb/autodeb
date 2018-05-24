// Package cli is responsible for parsing command line arguments and creating
// a server config.
//
// It isn't this package's responsibility to ensure that the configuration is
// valid.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
)

// Parse reads arguments and creates an autodeb server config
func Parse(args []string, writerOutput io.Writer) (*config.Config, error) {

	fs := flag.NewFlagSet("autodeb-server", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	var h, help bool
	fs.BoolVar(&help, "help", false, "print usage information")
	fs.BoolVar(&h, "h", false, "print usage information")

	var address string
	fs.StringVar(&address, "address", ":8071", "address to listen to")

	var templatesDirectory string
	fs.StringVar(&templatesDirectory, "templates-directory", "web/templates", "templates directory")

	var cacheTemplates bool
	fs.BoolVar(&cacheTemplates, "cache-templates", true, "whether or not to cache templates")

	var staticFilesDirectory string
	fs.StringVar(&staticFilesDirectory, "static-files-directory", "web/static", "static files directory")

	var dataDirectory string
	fs.StringVar(&dataDirectory, "data-directory", "data", "data directory")

	var databaseDriver string
	fs.StringVar(&databaseDriver, "database-driver", "sqlite3", "database driver")

	var databaseConnectionString string
	fs.StringVar(&databaseConnectionString, "database-connection-string", "database.sqlite", "database connection string")

	var logLevelString string
	fs.StringVar(&logLevelString, "log-level", "info", "info, warning or error")

	var oauthProvider string
	fs.StringVar(&oauthProvider, "oauth-provider", "gitlab", "oauth provider")

	var oauthBaseURL string
	fs.StringVar(&oauthBaseURL, "oauth-base-url", "https://salsa.debian.org", "oauth base url")

	var oauthClientID string
	fs.StringVar(&oauthClientID, "oauth-client-id", "", "oauth client id")

	var oauthClientSecret string
	fs.StringVar(&oauthClientSecret, "oauth-client-secret", "", "oauth client secret")

	var serverURLString string
	fs.StringVar(&serverURLString, "server-url", "http://localhost:8071", "public server url")

	var authentificationBackend string
	fs.StringVar(&authentificationBackend, "authentification-backend", "disabled", "selected authentification backend")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	for _, arg := range fs.Args() {
		if arg != "" {
			return nil, fmt.Errorf("unrecognized argument: %s", arg)
		}
	}

	if h || help {
		fs.SetOutput(writerOutput)
		fs.Usage()
		return nil, nil
	}

	serverURL, err := url.Parse(serverURLString)
	if err != nil {
		return nil, fmt.Errorf("invalid server url %s", serverURLString)
	}

	var logLevel log.Level
	switch logLevelString {
	case "info":
		logLevel = log.InfoLevel
	case "warning":
		logLevel = log.WarningLevel
	case "error":
		logLevel = log.ErrorLevel
	default:
		return nil, fmt.Errorf("unrecognized log level: %s", logLevelString)
	}

	cfg := &config.Config{
		DB: &config.DBConfig{
			Driver:           databaseDriver,
			ConnectionString: databaseConnectionString,
		},
		HTTP: &config.HTTPServerConfig{
			Address: address,
		},
		Auth: &config.AuthConfig{
			AuthentificationBackend: authentificationBackend,
			OAuth: &config.OAuthConfig{
				Provider:     oauthProvider,
				BaseURL:      oauthBaseURL,
				ClientID:     oauthClientID,
				ClientSecret: oauthClientSecret,
			},
		},
		ServerURL:             serverURL,
		DataDirectory:         dataDirectory,
		TemplatesDirectory:    templatesDirectory,
		StaticFilesDirectory:  staticFilesDirectory,
		TemplatesCacheEnabled: cacheTemplates,
		LogLevel:              logLevel,
	}

	return cfg, nil
}
