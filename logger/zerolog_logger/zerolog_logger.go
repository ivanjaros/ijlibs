package zerolog_logger

import (
	"github.com/ivanjaros/ijlibs/logger"
	"github.com/rs/zerolog"
)

// The logger is based on zerolog levels so we can simply convert the types and that's it.
func Convert(l logger.Level) (zerolog.Level, error) {
	return zerolog.Level(l), nil
}

func New(b logger.Logger) zerolog.Logger {
	l, _ := Convert(b.Level)
	return zerolog.New(b.Writer).With().Timestamp().Logger().Level(l)
}
