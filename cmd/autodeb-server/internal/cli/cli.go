// Package cli is responsible for parsing command line arguments and creating
// a server config.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server"
)

// Parse reads arguments and creates an autodeb server config
func Parse(args []string, writerOutput io.Writer) (*server.Config, error) {

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

	cfg := &server.Config{
		HTTP: server.HTTPServerConfig{
			Address: address,
		},
		DB: server.DBConfig{
			Driver:           databaseDriver,
			ConnectionString: databaseConnectionString,
		},
		DataDirectory:         dataDirectory,
		TemplatesDirectory:    templatesDirectory,
		StaticFilesDirectory:  staticFilesDirectory,
		TemplatesCacheEnabled: cacheTemplates,
		LogLevel:              logLevel,
	}

	return cfg, nil
}
