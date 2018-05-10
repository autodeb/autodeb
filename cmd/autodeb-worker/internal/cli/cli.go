// Package cli is responsible for parsing command line arguments and creating
// a worker config.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"salsa.debian.org/autodeb-team/autodeb/internal/worker"
)

// Parse reads arguments and creates an autodeb worker config
func Parse(args []string, writerOutput, writerError io.Writer) (*worker.Config, error) {

	fs := flag.NewFlagSet("autodeb-worker", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	var h, help bool
	fs.BoolVar(&help, "help", false, "print usage information")
	fs.BoolVar(&h, "h", false, "print usage information")

	var serverURL string
	fs.StringVar(&serverURL, "server-url", "", "URL of the autodeb server")

	var workingDirectory string
	fs.StringVar(&workingDirectory, "working-directory", "jobs", "working directory for jobs")

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

	cfg := &worker.Config{
		ServerURL:        serverURL,
		WorkingDirectory: workingDirectory,
		WriterOutput:     writerOutput,
		WriterError:      writerError,
	}

	return cfg, nil
}
