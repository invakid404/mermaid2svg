package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func SetupLogger(logLevel string, pretty bool) error {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DefaultContextLogger = &log.Logger

	parsedLogLevel, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level '%s': %w", logLevel, err)
	}
	log.Logger = log.Level(parsedLogLevel).With().Stack().Logger()

	if pretty {
		log.Logger = log.
			Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix}).
			With().
			Caller().
			Logger()
	}

	return nil
}
