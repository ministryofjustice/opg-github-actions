package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a configured slog.Logger instance
// that sets the log level and log handler.
//
// Will check env values (LOG_LEVEL, LOG_HANDLER) and
// replace the passed values with those if found
//
// By default, the level is set to Info and TextHandler
func New(lvl string, as string) (logger *slog.Logger) {
	var options = &slog.HandlerOptions{}

	if l := os.Getenv("LOG_LEVEL"); l != "" {
		lvl = l
	}
	if t := os.Getenv("LOG_HANDLER"); t != "" {
		as = t
	}

	switch lvl {
	case "ERROR", "error":
		options.Level = slog.LevelError
	case "WARN", "warn":
		options.Level = slog.LevelWarn
	case "INFO", "info":
		options.Level = slog.LevelInfo
	case "DEBUG", "debug":
		options.Level = slog.LevelDebug
	default:
		options.Level = slog.LevelInfo
	}

	as = strings.ToLower(as)
	if as == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, options))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, options))
	}
	return

}
