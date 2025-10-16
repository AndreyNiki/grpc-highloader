package logger

import (
	"context"
	"log/slog"
	"os"
)

const loggerKey = "ctxlogger"

type Logger struct {
	logger *slog.Logger
	file   *os.File
}

func NewLogFile(filepath string) (*Logger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	return &Logger{logger: logger}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) GetLogger() *slog.Logger {
	return l.logger
}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	defaultLogger := slog.Default()
	if ctx == nil {
		return defaultLogger
	}

	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}

	return defaultLogger
}
