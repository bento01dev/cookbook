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
		ctx := r.Context()
		h.ServeHTTP(w, r)
		slog.InfoContext(ctx, "request completed", "path", r.URL.Path, "method", r.Method, "elapsed_time_ms", time.Since(start).Milliseconds())
	})
}

func requestIdMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		request_id := r.Header.Get("request_id")
		if request_id == "" {
			ctx = context.WithValue(ctx, "request_id", uuid.New().String())
		} else {
			ctx = context.WithValue(ctx, "request_id", request_id)
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
