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

	initLog(stdout, getEnv("LOG_LEVEL"))
	slog.Info("log config set..")
	if err := startHttp(ctx, getEnv); err != nil {
		return fmt.Errorf("startup sequence for service failed..%w", err)
	}
	return nil
}

func initLog(w io.Writer, logLevelStr string) {
	var logLevel slog.Level = slog.LevelInfo
	if logLevelStr != "" {
		switch strings.ToLower(logLevelStr) {
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
	opts := slog.HandlerOptions{
		Level: logLevel,
	}

	logHandler := slog.NewJSONHandler(w, &opts)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
