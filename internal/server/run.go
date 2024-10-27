package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"log/slog"
)

func Run(ctx context.Context, stdout io.Writer, stderr io.Writer, args []string, getEnv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	initLog(stdout, getEnv)
	slog.Info("log config set..")
	if err := startHttp(ctx, getEnv); err != nil {
		return fmt.Errorf("startup sequence for service failed..%w", err)
	}
	return nil
}

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}
	return h.Handler.Handle(ctx, r)
}

func initLog(w io.Writer, getEnv func(string) string) {
	var logLevel slog.Level = slog.LevelInfo
	if v := getEnv("LOG_LEVEL"); v != "" {
		switch strings.ToLower(v) {
		case "debug":
			logLevel = slog.LevelDebug
		case "info":
			logLevel = slog.LevelInfo
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		}
	}
	defaultAttrs := []slog.Attr{
		slog.String("service", getEnv("SERVICE_NAME")),
		slog.String("env", getEnv("ENV")),
		slog.String("host", getEnv("HOST_IP")),
	}
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}

	baseHandler := slog.NewJSONHandler(w, &opts).WithAttrs(defaultAttrs)
	contextHandler := ContextHandler{Handler: baseHandler}

	logger := slog.New(&contextHandler)
	slog.SetDefault(logger)
}
