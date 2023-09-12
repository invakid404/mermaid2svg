package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func (api *API) NewRequestLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: api.log.Handler()})
}

type StructuredLogger struct {
	Logger slog.Handler
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	var logFields []slog.Attr
	logFields = append(logFields, slog.String("ts", time.Now().UTC().Format(time.RFC1123)))

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, slog.String("req_id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	handler := l.Logger.WithAttrs(append(logFields,
		slog.String("http_scheme", scheme),
		slog.String("http_proto", r.Proto),
		slog.String("http_method", r.Method),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))))

	entry := StructuredLoggerEntry{Logger: slog.New(handler)}

	entry.Logger.LogAttrs(context.Background(), slog.LevelInfo, "request started")

	return &entry
}

type StructuredLoggerEntry struct {
	Logger *slog.Logger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.LogAttrs(context.Background(), slog.LevelInfo, "request complete",
		slog.Int("resp_status", status),
		slog.Int("resp_byte_length", bytes),
		slog.Float64("resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0),
	)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, _ []byte) {
	l.Logger.LogAttrs(context.Background(), slog.LevelInfo, "", slog.String("panic", fmt.Sprintf("%+v", v)))
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

func GetLogEntry(r *http.Request) *slog.Logger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.With(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		for k, v := range fields {
			entry.Logger = entry.Logger.With(k, v)
		}
	}
}
