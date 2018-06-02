package log

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

// Level represents a logging level
// We don't use logrus's level because we want to implement UnmarshalText
// Also, we don't expose as many levels as logrus.
type Level int

// Logging levels
const (
	ErrorLevel Level = iota
	WarningLevel
	InfoLevel
)

func (level Level) String() string {
	switch level {
	case ErrorLevel:
		return "error"
	case WarningLevel:
		return "warning"
	case InfoLevel:
		return "info"
	}
	return "unknown"
}

// UnmarshalText will parse a log level
func (level *Level) UnmarshalText(text []byte) error {
	var err error

	s := string(text)

	switch s {
	case "error":
		*level = ErrorLevel
	case "warning":
		*level = WarningLevel
	case "info":
		*level = InfoLevel
	default:
		err = errors.Errorf("unrecognized log level %s", s)
	}

	return err
}
