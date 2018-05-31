// Package cli is responsible for parsing command line arguments and creating
// a worker config.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker"
)

// Parse reads arguments and creates an autodeb worker config
func Parse(args []string, writerOutput io.Writer) (*worker.Config, error) {

	fs := flag.NewFlagSet("autodeb-worker", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	var h, help bool
	fs.BoolVar(&help, "help", false, "print usage information")
	fs.BoolVar(&h, "h", false, "print usage information")

	var serverURL string
	fs.StringVar(&serverURL, "server-url", "", "URL of the autodeb server")

	var accessToken string
	fs.StringVar(&accessToken, "access-token", "", "API access token")

	var workingDirectory string
	fs.StringVar(&workingDirectory, "working-directory", "jobs", "working directory for jobs")

	var logLevelString string
	fs.StringVar(&logLevelString, "log-level", "info", "info, warning or error")

	var runnerCount int
	fs.IntVar(&runnerCount, "runner-count", 1, "number of job runners (concurrent jobs)")

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

	if serverURL == "" {
		return nil, fmt.Errorf("missing argument: server-url")
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

	cfg := &worker.Config{
		AccessToken:      accessToken,
		ServerURL:        serverURL,
		WorkingDirectory: workingDirectory,
		LogLevel:         logLevel,
		RunnerCount:      runnerCount,
	}

	return cfg, nil
}
