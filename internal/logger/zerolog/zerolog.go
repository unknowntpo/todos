package zerolog

import (
	"io"
	"os"

	"github.com/unknowntpo/todos/internal/logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type zerologWrapper struct {
	logger *zerolog.Logger
}

func New(out io.Writer) logger.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log := zerolog.New(out).With().Stack().Timestamp().Logger()

	return &zerologWrapper{logger: &log}
}

func (zw *zerologWrapper) PrintInfo(message string, properties map[string]interface{}) {
	zw.logger.WithLevel(zerolog.InfoLevel).
		Fields(properties).Msg(message)
}

func (zw *zerologWrapper) PrintError(err error, properties map[string]interface{}) {
	zw.logger.WithLevel(zerolog.ErrorLevel).Stack().
		Err(err).Fields(properties).Msg("")
}

func (zw *zerologWrapper) PrintFatal(err error, properties map[string]interface{}) {
	zw.logger.WithLevel(zerolog.FatalLevel).Stack().
		Err(err).Fields(properties).Msg("")
	os.Exit(1)
}
