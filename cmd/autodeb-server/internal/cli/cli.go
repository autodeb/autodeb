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

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
)

// Parse reads arguments and creates an autodeb server config
func Parse(args []string, fs filesystem.FS, writerOutput io.Writer) (*config.Config, error) {

	flagSet := flag.NewFlagSet("autodeb-server", flag.ContinueOnError)
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

	cfg, err := config.ParseConfig(configFile, fs)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
