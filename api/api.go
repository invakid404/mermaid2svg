package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/invakid404/mermaid2svg/api/routes"
	"github.com/invakid404/mermaid2svg/webdriver"
	"net/http"
)

type Options struct {
	Driver *webdriver.Driver
}

type API struct {
	server *http.Server
	router chi.Router
	driver *webdriver.Driver
}

func New(options Options) *API {
	api := &API{
		driver: options.Driver,
		router: chi.NewRouter(),
	}

	api.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), "driver", api.driver)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	})

	routes.Register(api.router)

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

	return nil
}

func (api *API) Stop() error {
	return api.server.Close()
}
