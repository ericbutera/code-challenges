package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	return NewWithLevel(slog.LevelInfo)
}

func NewWithLevel(level slog.Leveler) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
