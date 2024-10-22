package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/bento01dev/cookbook/internal/services"
)

func startHttp(ctx context.Context, getEnv func(string) string) error {
	var err error
	host := "127.0.0.1"
	if getEnv("HTTP_HOST") != "" {
		host = getEnv("HTTP_HOST")
	}
	port := "8080"
	if getEnv("HTTP_PORT") != "" {
		port = getEnv("HTTP_PORT")
	}

	// initialising and starting server..
	var rs recipeService
	if getEnv("MEMORY_REPO") != "" {
		rs, err = services.NewRecipeService(services.WithMemoryRepository())
		if err != nil {
			return err
		}
	}

	srv := NewServer(rs)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: srv,
	}
	go func() {
		slog.Info("starting cookbook service..", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			slog.Error("unable to start server", "err", err.Error())
		}
	}()

	// waiting for service interrupt based on ctx passed to gracefully shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		slog.Info("Shutting down service..", "addr", httpServer.Addr)
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown of server was not graceful", "err", err.Error())
		}
	}()
	wg.Wait()

	return nil
}

func NewServer(
	rs recipeService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, rs)
	//TODO: add middlewares
	var handler http.Handler = mux
	handler = NewRequestTimerMiddleware(handler)
	handler = NewRequestIdMiddleware(handler)
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	rs recipeService,
) {
	mux.Handle("GET /healthz", handleHealthz())

	mux.Handle("GET /recipe/{id}", handleGetRecipe(rs))
	mux.Handle("POST /recipe", handleCreateRecipe(rs))
}
