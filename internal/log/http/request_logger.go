package http

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// RequestLogger logs requests
type RequestLogger struct {
	logger log.Logger
}

//NewRequestLogger creates a new RequestLogger
func NewRequestLogger(logger log.Logger) *RequestLogger {
	requestLogger := &RequestLogger{
		logger: logger,
	}
	return requestLogger
}

// Error will log a request error
func (l *RequestLogger) Error(r *http.Request, err error) {
	line := fmt.Sprintf(
		"%s\t%s\t%+v",
		r.Method,
		r.RequestURI,
		err,
	)

	l.logger.Error(line)
}
