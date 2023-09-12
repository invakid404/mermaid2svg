package api

import (
	"errors"
	"fmt"
	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5"
	"github.com/invakid404/mermaid2svg/util/httputil"
	"github.com/invakid404/mermaid2svg/webdriver"
	"log/slog"
	"net/http"
)

type Options struct {
	Driver  *webdriver.Driver
	Log     *slog.Logger
	Checker health.Checker
}

type API struct {
	server  *http.Server
	router  chi.Router
	driver  *webdriver.Driver
	log     *slog.Logger
	checker health.Checker
}

func New(options Options) *API {
	api := &API{
		driver:  options.Driver,
		log:     options.Log,
		router:  chi.NewRouter(),
		checker: options.Checker,
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
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("failed to start api: %w", err))
		}
	}()

	api.log.Info(fmt.Sprintf("api listening on %s", api.server.Addr))

	return nil
}

func (api *API) Stop() error {
	if err := httputil.ShutdownGracefully(api.server); err != nil {
		return fmt.Errorf("failed to shutdown api server: %w", err)
	}

	return nil
}
