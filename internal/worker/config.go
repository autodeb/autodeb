package worker

import (
	"io"
)

// Config contains configuration for Worker
type Config struct {
	ServerURL    string
	WriterOutput io.Writer
	WriterError  io.Writer
}
