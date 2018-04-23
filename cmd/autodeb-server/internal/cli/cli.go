// Package cli is responsible for parsing command line arguments and creating
// a server instance.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"salsa.debian.org/autodeb-team/autodeb/internal/http"
	"salsa.debian.org/autodeb-team/autodeb/internal/logo"
	"salsa.debian.org/autodeb-team/autodeb/internal/server"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// Run reads arguments and creates an autodeb server
func Run(args []string, writerOutput, writerError io.Writer) (*server.Server, error) {

	fs := flag.NewFlagSet("autodeb-server", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	var h, help bool
	fs.BoolVar(&help, "help", false, "print usage information")
	fs.BoolVar(&h, "h", false, "print usage information")

	var address string
	fs.StringVar(&address, "address", "0.0.0.0", "address to listen to")

	var port int
	fs.IntVar(&port, "port", 8080, "port to listen to")

	var templatesDirectory string
	fs.StringVar(&templatesDirectory, "templates", "web/templates", "templates directory")

	var dataDirectory string
	fs.StringVar(&dataDirectory, "data", "data", "data directory")

	var databaseDriver string
	fs.StringVar(&databaseDriver, "database-driver", "sqlite3", "database driver")

	var databaseConnectionString string
	fs.StringVar(&databaseConnectionString, "database-connection-string", "database.sqlite", "database connection string")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() > 0 {
		return nil, fmt.Errorf("unrecognized argument: %s", fs.Arg(0))
	}

	if h || help {
		fs.SetOutput(writerOutput)
		fs.Usage()
		return nil, nil
	}

	fmt.Fprintln(writerOutput, logo.Logo)

	fmt.Fprintf(writerOutput, "Starting autodeb API on %s:%d.\n", address, port)

	cfg := &server.Config{
		App: &app.Config{
			TemplatesDirectory: templatesDirectory,
			DataDirectory:      dataDirectory,
		},
		HTTP: &http.ServerConfig{
			Address: address,
			Port:    port,
		},
		Database: &database.Config{
			Driver:           databaseDriver,
			ConnectionString: databaseConnectionString,
		},
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
