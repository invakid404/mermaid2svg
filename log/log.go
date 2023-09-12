package log

import (
	"fmt"
	"log/slog"
	"os"
)

func SetupLogger(logLevel string, pretty bool) error {
	parsedLogLevel := slog.LevelInfo
	if logLevel != "" {
		if err := parsedLogLevel.UnmarshalText([]byte(logLevel)); err != nil {
			return fmt.Errorf("failed to parse log level: %w", err)
		}
	}

	handlerOptions := &slog.HandlerOptions{
		Level: parsedLogLevel,
	}

	var logger *slog.Logger
	if pretty {
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	}

	slog.SetDefault(logger)

	return nil
}
