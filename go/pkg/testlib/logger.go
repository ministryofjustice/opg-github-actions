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

func Testlogger(opts *slog.HandlerOptions) {
	var handler *slog.TextHandler
	if opts == nil {
		opts = VerboseOptions
	}
	handler = slog.NewTextHandler(stdout, opts)
	slog.SetDefault(slog.New(handler))
}
