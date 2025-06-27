package logger

import (
	"context"
	"go-fhir-demo/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm/logger"
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

	// Set log format with enhanced coloring
	if format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     true,
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			DisableColors:   false,
			PadLevelText:    true,
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

// GormLogger implements gormLogger.Interface and injects trace/span IDs into SQL logs.
type GormLogger struct {
	logger        *logrus.Logger
	logLevel      logger.LogLevel
	slowThreshold time.Duration
	ctx           context.Context
}

// GetGormLogger returns a GormLogger with the provided context.
func GetGormLogger(ctx context.Context) logger.Interface {
	return &GormLogger{
		logger:        GetLogger(),
		logLevel:      logger.Info,
		slowThreshold: 200 * time.Millisecond,
		ctx:           ctx,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		entry := l.entryWithTrace(ctx)
		entry.Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		entry := l.entryWithTrace(ctx)
		entry.Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		entry := l.entryWithTrace(ctx)
		entry.Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	entry := l.entryWithTrace(ctx)

	// Clean up SQL query - remove extra spaces and normalize whitespace
	cleanSQL := strings.TrimSpace(strings.ReplaceAll(sql, "\n", " "))
	cleanSQL = strings.TrimSpace(strings.ReplaceAll(cleanSQL, "\\", " "))
	cleanSQL = strings.Join(strings.Fields(cleanSQL), " ")

	// Redact PII fields from SQL before logging
	cleanSQL = utils.RedactPIIFromSQL(cleanSQL)

	// Extract trace/span IDs from context
	traceID := ""
	spanID := ""
	span := trace.SpanFromContext(ctx)
	if span != nil {
		sc := span.SpanContext()
		if sc.HasTraceID() {
			traceID = sc.TraceID().String()
		}
		if sc.HasSpanID() {
			spanID = sc.SpanID().String()
		}
	}

	// Enhanced formatting with colors and better spacing
	var logMsg string

	// Color codes for better visibility
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[34m"
		colorGreen  = "\033[32m"
		colorCyan   = "\033[36m"
		colorBold   = "\033[1m"
	)

	// Compose trace/span info line
	traceSpanLine := colorCyan + "║ " + colorBold + "trace_id=" + colorReset + colorBlue + traceID + colorReset +
		" " + colorBold + "span_id=" + colorReset + colorBlue + spanID + colorReset + "\n"

	if err != nil {
		logMsg = "\n" + colorBold + colorRed + "╔═══ SQL ERROR ═══════════════════════════════════════════════════════════════════╗" + colorReset + "\n" +
			colorRed + "║ " + colorBold + "QUERY: " + colorReset + colorRed + cleanSQL + colorReset + "\n" +
			colorRed + "║ " + colorBold + "ERROR: " + colorReset + colorRed + err.Error() + colorReset + "\n" +
			colorRed + "║ " + colorBold + "TIME:  " + colorReset + colorRed + elapsed.String() + colorReset + "\n" +
			colorRed + "║ " + colorBold + "ROWS:  " + colorReset + colorRed + "%d" + colorReset + "\n" +
			traceSpanLine +
			colorBold + colorRed + "╚═════════════════════════════════════════════════════════════════════════════════╝" + colorReset + "\n"
		logMsg = strings.Replace(logMsg, "%d", formatRows(rows), 1)
		entry.Error(logMsg)
	} else if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		logMsg = "\n" + colorBold + colorYellow + "╔═══ SLOW SQL QUERY ══════════════════════════════════════════════════════════════╗" + colorReset + "\n" +
			colorYellow + "║ " + colorBold + "QUERY: " + colorReset + colorYellow + cleanSQL + colorReset + "\n" +
			colorYellow + "║ " + colorBold + "TIME:  " + colorReset + colorYellow + elapsed.String() + colorReset + "\n" +
			colorYellow + "║ " + colorBold + "ROWS:  " + colorReset + colorYellow + "%d" + colorReset + "\n" +
			traceSpanLine +
			colorBold + colorYellow + "╚═════════════════════════════════════════════════════════════════════════════════╝" + colorReset + "\n"
		logMsg = strings.Replace(logMsg, "%d", formatRows(rows), 1)
		entry.Warn(logMsg)
	} else {
		logMsg = "\n" + colorBold + colorCyan + "╔═══ SQL QUERY ═══════════════════════════════════════════════════════════════════╗" + colorReset + "\n" +
			colorCyan + "║ " + colorBold + "QUERY: " + colorReset + colorGreen + cleanSQL + colorReset + "\n" +
			colorCyan + "║ " + colorBold + "TIME:  " + colorReset + colorBlue + elapsed.String() + colorReset + "\n" +
			colorCyan + "║ " + colorBold + "ROWS:  " + colorReset + colorBlue + "%d" + colorReset + "\n" +
			traceSpanLine +
			colorBold + colorCyan + "╚═════════════════════════════════════════════════════════════════════════════════╝" + colorReset + "\n"
		logMsg = strings.Replace(logMsg, "%d", formatRows(rows), 1)
		entry.Info(logMsg)
	}

	// Send SQL structure as event to current span if present
	if span != nil && span.IsRecording() {
		span.AddEvent("SQL",
			trace.WithAttributes(
				attribute.String("db.statement", cleanSQL),
				attribute.Float64("db.elapsed_ms", float64(elapsed.Microseconds())/1000.0),
				attribute.Int64("db.rows", rows),
			),
		)
	}
}

// formatRows returns the string representation of the rows value.
func formatRows(rows int64) string {
	return strconv.FormatInt(rows, 10)
}

func (l *GormLogger) entryWithTrace(_ context.Context) *logrus.Entry {
	// Do NOT add trace_id/span_id fields for SQL logs to avoid them appearing outside the box
	return l.logger.WithFields(logrus.Fields{})
}

type TraceID interface {
	String() string
}
type SpanID interface {
	String() string
}
