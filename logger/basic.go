package logger

import (
	"io"
	"os"
)

// the simplest wrapper containing data that can be used to create instance
// of any logging library.
type Logger struct {
	Writer io.Writer
	Level  Level
}

func New() Logger {
	return Logger{Writer: os.Stdout, Level: ErrorLevel}
}
