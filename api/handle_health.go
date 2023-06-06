package api

import (
	"github.com/alexliesenfeld/health"
	"net/http"
)

func (api *API) registerHealth() {
	api.router.Get("/livez", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte("pong"))
	})

	api.router.Get("/readyz", health.NewHandler(api.checker))
}
