package log

import (
	"fmt"
	"strings"
)

type prefixLogger struct {
	logger Logger
	prefix string
}

func newPrefixLogger(logger Logger, prefix string) Logger {
	prefix = strings.TrimSpace(prefix)
	prefix = fmt.Sprintf("[%s] ", prefix)

	prefixLogger := &prefixLogger{
		logger: logger,
		prefix: prefix,
	}

	return prefixLogger
}

// PrefixLogger returns a sub-logger that uses a prefix
func (l *prefixLogger) PrefixLogger(prefix string) Logger {
	return newPrefixLogger(l, prefix)
}

// SetLevel sets the logging level
func (l *prefixLogger) SetLevel(level Level) {
	l.logger.SetLevel(level)
}

// Error logs an error message
func (l *prefixLogger) Error(args ...interface{}) {
	args = append([]interface{}{l.prefix}, args)
	l.logger.Error(args...)
}

// Errorf logs an error message with the given format
func (l *prefixLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(l.prefix+format, args...)
}

// Warning logs a warning message
func (l *prefixLogger) Warning(args ...interface{}) {
	args = append([]interface{}{l.prefix}, args)
	l.logger.Warning(args...)
}

// Warningf logs a warning message with the given format
func (l *prefixLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(l.prefix+format, args...)
}

// Info logs a info message
func (l *prefixLogger) Info(args ...interface{}) {
	args = append([]interface{}{l.prefix}, args)
	l.logger.Info(args...)
}

// Infof logs a info message with the given format
func (l *prefixLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(l.prefix+format, args...)
}
