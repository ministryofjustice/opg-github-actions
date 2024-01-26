package testlib

import (
	"io"
	"log/slog"
	"os"
)

var (
	DefaultOptions = &slog.HandlerOptions{Level: slog.LevelError}
	InfoOptions    = &slog.HandlerOptions{Level: slog.LevelInfo}
	VerboseOptions = &slog.HandlerOptions{Level: slog.LevelDebug}
	ignore         = io.Discard
	stdout         = os.Stdout
)
var (
	logLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

func Testlogger(opts *slog.HandlerOptions) {
	var (
		handler *slog.TextHandler
		level   string = os.Getenv("LOG_LEVEL") // "error"
		to      string = os.Getenv("LOG_TO")
	)

	if opts == nil && level != "" {
		opts = &slog.HandlerOptions{Level: logLevels[level]}
	} else if opts == nil {
		opts = DefaultOptions
	}

	if to == "stdout" {
		handler = slog.NewTextHandler(stdout, opts)
	} else {
		handler = slog.NewTextHandler(ignore, opts)
	}

	slog.SetDefault(slog.New(handler))
}
