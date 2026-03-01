package logger

import (
	"log/slog"
	"os"
	"github.com/muchirisworld/terminal/internal/config"
)

// New creates a new slog.Logger.
func New(cfg *config.Config) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		AddSource: true,
	}


	switch cfg.LogLevel {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	switch cfg.AppEnv {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
