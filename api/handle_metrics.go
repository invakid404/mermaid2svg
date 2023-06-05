package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (api *API) registerMetrics() {
	api.router.Get("/metrics", promhttp.Handler().ServeHTTP)
}
