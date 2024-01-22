package testlib

import (
	"log/slog"
	"os"
)

var (
	DefaultOptions = &slog.HandlerOptions{Level: slog.LevelError}
	InfoOptions    = &slog.HandlerOptions{Level: slog.LevelInfo}
	VerboseOptions = &slog.HandlerOptions{Level: slog.LevelDebug}
	//destination    = io.Discard
	destination = os.Stdout
)

func Testlogger(opts *slog.HandlerOptions) {
	if opts == nil {
		opts = DefaultOptions
	}
	handler := slog.NewTextHandler(destination, opts)
	slog.SetDefault(slog.New(handler))
}
