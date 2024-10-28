package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/bento01dev/cookbook/internal/config"
	"github.com/bento01dev/cookbook/internal/services"
	"github.com/bento01dev/cookbook/internal/stats"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startHttp(ctx context.Context, getEnv func(string) string) error {
	var (
		err  error
		conf config.Config
	)
	conf, err = config.NewConfig(getEnv)
	if err != nil {
		return fmt.Errorf("unable to generate config:%w", err)
	}

	statsCollection := stats.Stats(getEnv)

	// initialising and starting server..
	var rs recipeService
	if getEnv("MEMORY_REPO") != "" {
		rs, err = services.NewRecipeService(services.WithMemoryRepository())
		if err != nil {
			return err
		}
	}

	srv := NewServer(rs, statsCollection, conf)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(conf.Host, conf.Port),
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
	statsCollection *stats.StatsCollection,
	conf config.Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, rs, statsCollection, conf)

	var handler http.Handler = mux
	handler = requestTimerMiddleware(handler)
	handler = requestIdMiddleware(handler)
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	rs recipeService,
	statsCollection *stats.StatsCollection,
	conf config.Config,
) {
	mux.Handle("GET /healthz", handleHealthz())

	mux.Handle("GET /recipe/{id}", timeoutMiddleware(handleGetRecipe(rs, statsCollection), conf.GetRecipeTimeout))
	mux.Handle("POST /recipe", timeoutMiddleware(handleCreateRecipe(rs, statsCollection), conf.CreateRecipeTimeout))

	mux.Handle("/metrics", promhttp.Handler())
}
