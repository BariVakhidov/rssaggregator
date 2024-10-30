package logger

import (
	"log/slog"
	"os"
)

const (
	modeLocal = "local"
	modeProd  = "prod"
)

type Logger struct {
	Log *slog.Logger
}

func New(loggerMode string) *Logger {
	var log *slog.Logger

	switch loggerMode {
	case modeLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case modeProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return &Logger{Log: log}
}
