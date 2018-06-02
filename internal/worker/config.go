package worker

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// Config contains configuration for Worker
type Config struct {
	ServerURL        string    `toml:"server_url"`
	AccessToken      string    `toml:"access_token"`
	WorkingDirectory string    `toml:"working_directory"`
	LogLevel         log.Level `toml:"log_level"`
	RunnerCount      int       `toml:"runner_count"`
}
