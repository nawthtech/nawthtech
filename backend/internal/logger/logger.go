package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type DefaultLogger struct {
	logger *slog.Logger
}

var (
	Stdout *slog.Logger
	Stderr *slog.Logger

	globalLogger Logger
)

func Init(env string) {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if env == "development" {
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
		opts.AddSource = true
		Stdout = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	globalLogger = &DefaultLogger{logger: Stdout}
}

func (l *DefaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *DefaultLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *DefaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *DefaultLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *DefaultLogger) With(args ...any) Logger {
	return &DefaultLogger{logger: l.logger.With(args...)}
}

func Debug(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Debug(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Info(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Warn(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Error(ctx, msg, args...)
}

// Helpers
func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "")
	}
	return slog.String("error", err.Error())
}

func DurationAttr(duration time.Duration) slog.Attr {
	return slog.Duration("duration", duration)
}

func TimestampAttr() slog.Attr {
	return slog.String("timestamp", time.Now().Format(time.RFC3339))
}

func RequestAttr(method, path string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
	)
}

func formatMemory(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string(rune(float64(bytes)/float64(div))) + " " + string("KMGTPE"[exp]) + "B"
}

func GetGlobalLogger() Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger
}

func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}