package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger interface {
	Panic(err error, message string)
	Fatal(err error, message string)
	Errorf(err error, format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
}

type customLog struct {
	logger zerolog.Logger
}

func NewLogger() Logger {
	//nolint:reassign
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return customLog{logger: log.Logger}
}

func (log customLog) Panic(err error, message string) {
	log.logger.Panic().Stack().Err(err).Msg(message)
}

func (log customLog) Fatal(err error, message string) {
	log.logger.Fatal().Stack().Err(err).Msg(message)
}

func (log customLog) Errorf(err error, format string, args ...interface{}) {
	log.logger.Error().Stack().Err(err).Msgf(format, args...)
}

func (log customLog) Warnf(format string, args ...interface{}) {
	log.logger.Warn().Msgf(format, args...)
}

func (log customLog) Infof(format string, args ...interface{}) {
	log.logger.Info().Msgf(format, args...)
}

func (log customLog) Debugf(format string, args ...interface{}) {
	log.logger.Debug().Msgf(format, args...)
}

func (log customLog) Tracef(format string, args ...interface{}) {
	log.logger.Trace().Msgf(format, args...)
}
