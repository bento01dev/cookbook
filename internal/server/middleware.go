package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func requestTimerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		if request_id, ok := r.Context().Value("request_id").(string); ok {
			slog.Info("request completed", "path", r.URL.Path, "method", r.Method, "request_id", request_id, "elapsed_time_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("request completed", "path", r.URL.Path, "method", r.Method, "elapsed_time_ms", time.Since(start).Milliseconds())
	})
}

func requestIdMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		request_id := ctx.Value("request_id")
		if request_id == nil {
			ctx = context.WithValue(ctx, "request_id", uuid.New().String())
		}
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func timeoutMiddleware(h http.Handler, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
