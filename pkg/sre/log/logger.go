package log

import (
	"log/slog"
	"os"
)

func ToSLogLevel(logLevel string) slog.Level {
	slogLevel := slog.LevelError
	switch logLevel {
	case "debug":
		slogLevel = slog.LevelDebug
	case "error":
		slogLevel = slog.LevelError
	case "warn":
		slogLevel = slog.LevelWarn
	case "info":
		slogLevel = slog.LevelInfo
	default:
		slog.With("target_level", logLevel).Error("level not handled, fallback on error")
	}
	return slogLevel
}

func SetupLogger(logLevel slog.Level) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel,
		ReplaceAttr: nil,
	})))
}
