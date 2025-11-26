package logger

import (
	"os"

	"log/slog"
)

var (
	stdoutHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	//enable source
	stdoutHandlerWithSource = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})

	stderrHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	// enable source
	stderrHandlerWithSource = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})

	// sends logs to stdout
	Stdout = slog.New(stdoutHandler)
	// sends logs to stdout with source info
	StdoutWithSource = slog.New(stdoutHandlerWithSource)

	// sends logs to stderr
	Stderr = slog.New(stderrHandler)
	// sends logs to stderr with source info
	StderrWithSource = slog.New(stderrHandlerWithSource)
)

// ErrAttr دالة مساعدة لإرجاع سمة الخطأ
func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

// ErrorsAttr دالة مساعدة لإرجاع سمة الأخطاء المتعددة
func ErrorsAttr(errors ...error) slog.Attr {
	return slog.Any("errors", errors)
}
