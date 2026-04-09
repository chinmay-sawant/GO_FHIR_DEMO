// Package logger provides structured logging with trace ID support.
package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var appLogger = logrus.New()

// Initialize sets up the logger with the given configuration.
func Initialize(level, format, logFile string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	appLogger.SetLevel(logLevel)

	if format == "json" {
		appLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     true,
		})
	} else {
		appLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			DisableColors:   false,
			PadLevelText:    true,
		})
	}

	if logFile != "" {
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		appLogger.SetOutput(io.MultiWriter(os.Stdout, appendFileWriter{path: logFile}))
	} else {
		appLogger.SetOutput(os.Stdout)
	}

	return nil
}

// Close is retained for backward compatibility.
func Close() error {
	return nil
}

// GetLogger returns the configured logger instance.
func GetLogger() *logrus.Logger {
	return appLogger
}

// WithContext returns a log entry with context fields for tracing.
func WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	span := trace.SpanFromContext(ctx)
	if span != nil {
		sc := span.SpanContext()
		if sc.HasTraceID() {
			fields["trace_id"] = sc.TraceID().String()
		}
		if sc.HasSpanID() {
			fields["span_id"] = sc.SpanID().String()
		}
	}

	return GetLogger().WithFields(fields)
}

type appendFileWriter struct {
	path string
}

func (w appendFileWriter) Write(p []byte) (int, error) {
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()
	return file.Write(p)
}
