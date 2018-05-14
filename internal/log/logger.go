// Package log provides a logger. The logger currently wraps sirupsen/logrus's
// Logger but it could be easily replaced.
package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logger is used to log error, warning and info messages
type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Warning(...interface{})
	Warningf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	PrefixLogger(string) Logger
	SetLevel(Level)
}

// Level represents a logging level
type Level = logrus.Level

// Logging levels
const (
	ErrorLevel   Level = logrus.ErrorLevel
	WarningLevel Level = logrus.WarnLevel
	InfoLevel    Level = logrus.InfoLevel
)

type logger struct {
	logger *logrus.Logger
}

// New creates a new logger
func New(output io.Writer) Logger {
	logrusLogger := logrus.New()
	logrusLogger.Out = output

	l := &logger{
		logger: logrusLogger,
	}

	return l
}

// PrefixLogger returns a sub-logger that uses a prefix
func (l *logger) PrefixLogger(prefix string) Logger {
	return newPrefixLogger(l, prefix)
}

// SetLevel sets the logging level
func (l *logger) SetLevel(level Level) {
	l.logger.SetLevel(level)
}

// Error logs an error message
func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf logs an error message with the given format
func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Warning logs a warning message
func (l *logger) Warning(args ...interface{}) {
	l.logger.Warning(args...)
}

// Warningf logs a warning message with the given format
func (l *logger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}

// Info logs a info message
func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logs a info message with the given format
func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
