package worker

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// Config contains configuration for Worker
type Config struct {
	ServerURL        string
	AccessToken      string
	WorkingDirectory string
	LogLevel         log.Level
	RunnerCount      int
}
