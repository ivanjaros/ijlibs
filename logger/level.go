// This is a close copy of the zerolog levels.
// We are using this centralised version so that we can use different logging backends
// and not be tied to a single logging library.
package logger

import (
	"fmt"
)

// Level defines log levels.
type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case NoLevel:
		return ""
	}
	return ""
}

// LevelFieldMarshalFunc allows customization of global level field marshaling
func LevelFieldMarshalFunc(l Level) string {
	return l.String()
}

// ParseLevel converts a level string into a zerolog Level value.
// returns an error if the input string does not match known values.
func ParseLevel(levelStr string) (Level, error) {
	switch levelStr {
	case LevelFieldMarshalFunc(TraceLevel):
		return TraceLevel, nil
	case LevelFieldMarshalFunc(DebugLevel):
		return DebugLevel, nil
	case LevelFieldMarshalFunc(InfoLevel):
		return InfoLevel, nil
	case LevelFieldMarshalFunc(WarnLevel), "warning": // warning is quite common term, we'll specifically support it
		return WarnLevel, nil
	case LevelFieldMarshalFunc(ErrorLevel):
		return ErrorLevel, nil
	case LevelFieldMarshalFunc(FatalLevel):
		return FatalLevel, nil
	case LevelFieldMarshalFunc(PanicLevel):
		return PanicLevel, nil
	case LevelFieldMarshalFunc(NoLevel):
		return NoLevel, nil
	}
	return NoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to NoLevel", levelStr)
}
