package logging

import (
	"log/slog"
	"os"
)

func LogLevelStringToSlogLevel(logLevel string) slog.Level {
	var l slog.Level
	switch logLevel {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelError
	}
	return l
}

func NewSlogLogger(logLevel string, pretty bool) *slog.Logger {
	l := LogLevelStringToSlogLevel(logLevel)
	if !pretty {
		return NewJSONSlogLogger(l)
	}
	return NewTextSlogLogger(l)
}
func NewJSONSlogLogger(logLevel slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
}
func NewTextSlogLogger(logLevel slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
}
