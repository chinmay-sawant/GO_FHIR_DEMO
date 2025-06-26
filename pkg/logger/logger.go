package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var Logger *logrus.Logger

// Initialize sets up the logger with the given configuration
func Initialize(level, format, logFile string) error {
	Logger = logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Logger.SetLevel(logLevel)

	// Set log format
	if format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set up log file output
	if logFile != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		// Write to both file and stdout
		multiWriter := io.MultiWriter(os.Stdout, file)
		Logger.SetOutput(multiWriter)
	}

	return nil
}

// GetLogger returns the configured logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		Logger = logrus.New()
	}
	return Logger
}

// Info logs an info message
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// WithContext returns a log entry with context fields for tracing (Jaeger/OpenTelemetry)
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

type SpanContext interface {
	TraceID() TraceID
	SpanID() SpanID
}
type TraceID interface {
	String() string
}
type SpanID interface {
	String() string
}
