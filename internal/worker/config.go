package worker

import (
	"io"
)

// Config contains configuration for Worker
type Config struct {
	ServerAddress string
	ServerPort    int
	WriterOutput  io.Writer
	WriterError   io.Writer
}
