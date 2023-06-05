package main

import (
	_ "embed"
	"fmt"
	"github.com/invakid404/mermaid2svg/api"
	logconfig "github.com/invakid404/mermaid2svg/log"
	"github.com/invakid404/mermaid2svg/webdriver"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func run() error {
	logLevel := os.Getenv("LOG_LEVEL")
	dev := os.Getenv("DEV") == "true"

	if logLevel == "" {
		logLevel = "info"
	}

	if err := logconfig.SetupLogger(logLevel, dev); err != nil {
		return fmt.Errorf("failed to set up logging: %w", err)
	}

	driver, err := webdriver.New()
	if err != nil {
		return fmt.Errorf("failed to init webdriver: %w", err)
	}

	defer func(driver *webdriver.Driver) {
		if err := driver.Stop(); err != nil {
			log.Error().Err(err).Msg("failed to stop web driver")
		}
	}(driver)

	if err = driver.Start(); err != nil {
		return fmt.Errorf("failed to start driver: %w", err)
	}

	app := api.New(api.Options{
		Driver: driver,
		Log:    log.With().Str("component", "api").Logger(),
	})

	defer func(app *api.API) {
		err := app.Stop()
		if err != nil {
			log.Error().Err(err).Msg("failed to stop api")
		}
	}(app)

	if err = app.Start(); err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}
