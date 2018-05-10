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

	var serverAddress string
	fs.StringVar(&serverAddress, "server-address", "", "address of the autodeb server")

	var serverPort int
	fs.IntVar(&serverPort, "server-port", 0, "port of the autodeb server")

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

	if serverAddress == "" {
		return nil, fmt.Errorf("missing argument: server-address")
	}

	if serverPort == 0 {
		return nil, fmt.Errorf("missing argument: server-port")
	}

	cfg := &worker.Config{
		ServerAddress: serverAddress,
		ServerPort:    serverPort,
		WriterOutput:  writerOutput,
		WriterError:   writerError,
	}

	return cfg, nil
}
