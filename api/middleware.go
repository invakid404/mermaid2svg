package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (api *API) RequestLogger(next http.Handler) http.Handler {
	return chi.Chain(
		hlog.NewHandler(api.log),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			l := log.Ctx(r.Context())
			l.Info().
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("HTTP request handled")
		}),
		hlog.URLHandler("url"),
		hlog.MethodHandler("method"),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
		hlog.RequestIDHandler("request_id", middleware.RequestIDHeader),
	).Handler(next)
}
