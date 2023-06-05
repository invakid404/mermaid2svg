package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (api *API) registerRoutes() {
	api.router.Use(
		middleware.RequestID,
		api.RequestLogger,
		cors.New(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
			},
		).Handler,
	)

	api.router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(middleware.AllowContentEncoding("application/json"))

		render := &renderAPI{api}
		render.Register(apiRouter)
	})
}
