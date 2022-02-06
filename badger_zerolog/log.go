package badger_zerolog

import (
	"github.com/rs/zerolog"
	"os"
)

func New() *logger {
	return &logger{l: zerolog.New(os.Stdout)}
}

func Wrap(l zerolog.Logger) *logger {
	return &logger{l: l}
}

type logger struct {
	l zerolog.Logger
}

func (l *logger) Errorf(msg string, args ...interface{}) {
	l.l.Error().Msgf(msg, args...)
}

func (l *logger) Warningf(msg string, args ...interface{}) {
	l.l.Warn().Msgf(msg, args...)
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.l.Info().Msgf(msg, args...)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	l.l.Debug().Msgf(msg, args...)
}
