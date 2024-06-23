package logger

import (
	"log/slog"
	"os"
)

func Logger() *slog.Logger {
	logHandler := slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}).WithAttrs([]slog.Attr{slog.String("version", "v1")})
	logger := slog.New(logHandler)
	return logger
}
