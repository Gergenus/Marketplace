package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(level string) *slog.Logger {
	switch level {
	case "debug":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "info":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("incorrect logger level")
	}
}
