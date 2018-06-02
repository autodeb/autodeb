package worker

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// ParseConfig parses a configuration file to create a worker config
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

	// Create the config, with defaults
	config := &Config{
		WorkingDirectory: "jobs",
		LogLevel:         log.InfoLevel,
		RunnerCount:      1,
		ServerURL:        "http://localhost:8071",
	}

	if metadata, err := toml.DecodeReader(bytes.NewReader(b), &config); err != nil {
		return nil, errors.WithMessage(err, "could not decode configuration file")
	} else if keys := metadata.Undecoded(); len(keys) > 0 {
		return nil, errors.Errorf("unrecognized configuration key: %s", keys[0].String())
	}

	return config, nil
}
