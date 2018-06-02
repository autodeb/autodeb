// Package cli is responsible for parsing command line arguments and creating
// a worker config.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker"
)

// Parse reads arguments and returns the worker configuration file path
func Parse(args []string, fs filesystem.FS, writerOutput io.Writer) (*worker.Config, error) {

	flagSet := flag.NewFlagSet("autodeb-worker", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	var h, help bool
	flagSet.BoolVar(&help, "help", false, "print usage information")
	flagSet.BoolVar(&h, "h", false, "print usage information")

	var configFile string
	flagSet.StringVar(&configFile, "config", "autodeb-worker.cfg", "path to configuration file")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	for _, arg := range flagSet.Args() {
		if arg != "" {
			return nil, fmt.Errorf("unrecognized argument: %s", arg)
		}
	}

	if h || help {
		flagSet.SetOutput(writerOutput)
		flagSet.Usage()
		return nil, nil
	}

	config, err := worker.ParseConfig(configFile, fs)
	if err != nil {
		return nil, err
	}

	return config, nil
}
