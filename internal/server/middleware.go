package server

import (
	"log/slog"
	"net/http"
	"time"
)

func NewRequestTimerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		slog.Info("request completed", "path", r.URL.Path, "method", r.Method, "elapsed_time_ms", time.Since(start).Milliseconds())
	})
}
