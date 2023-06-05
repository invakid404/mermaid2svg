package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/invakid404/mermaid2svg/util/httputil"
	"github.com/invakid404/mermaid2svg/webdriver"
	"github.com/rs/zerolog"
	"net/http"
)

type Options struct {
	Driver *webdriver.Driver
	Log    zerolog.Logger
}

type API struct {
	server *http.Server
	router chi.Router
	driver *webdriver.Driver
	log    zerolog.Logger
}

func New(options Options) *API {
	api := &API{
		driver: options.Driver,
		log:    options.Log,
		router: chi.NewRouter(),
	}

	api.registerRoutes()

	api.server = &http.Server{
		Addr:    ":8080",
		Handler: api.router,
	}

	return api
}

func (api *API) Start() error {
	go func() {
		err := api.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start api: %w", err))
		}
	}()

	api.log.Info().Msgf("API listening on %s", api.server.Addr)

	return nil
}

func (api *API) Stop() error {
	if err := httputil.ShutdownGracefully(api.server); err != nil {
		return fmt.Errorf("failed to shutdown api server: %w", err)
	}

	return nil
}
