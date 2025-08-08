package logger

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level represents the logging level
type Level = zapcore.Level

const (
	DEBUG Level = zapcore.DebugLevel
	INFO  Level = zapcore.InfoLevel
	WARN  Level = zapcore.WarnLevel
	ERROR Level = zapcore.ErrorLevel
	FATAL Level = zapcore.FatalLevel

	// Common configuration values
	timestampKey = "timestamp"
)

// Fields represents structured logging fields
type Fields = map[string]interface{}

// Logger interface for structured logging
type Logger interface {
	Debug(ctx context.Context, message string, fields Fields)
	Info(ctx context.Context, message string, fields Fields)
	Warn(ctx context.Context, message string, fields Fields)
	Error(ctx context.Context, message string, err error, fields Fields)
	Fatal(ctx context.Context, message string, err error, fields Fields)
	WithContext(ctx context.Context) Logger
	WithFields(fields Fields) Logger
}

// ZapLogger implements Logger interface using Zap
type ZapLogger struct {
	logger *zap.Logger
	fields Fields
}

// New creates a new ZapLogger with the specified level
func New(level Level) Logger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.TimeKey = timestampKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.DisableCaller = false
	config.DisableStacktrace = false

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		logger: logger,
		fields: make(Fields),
	}
}

// NewProduction creates a production-ready logger
func NewProduction() Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = timestampKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.DisableCaller = false
	config.DisableStacktrace = false

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		logger: logger,
		fields: make(Fields),
	}
}

// NewDevelopment creates a development logger with human-readable output
func NewDevelopment() Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = timestampKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.DisableCaller = false
	config.DisableStacktrace = false

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		logger: logger,
		fields: make(Fields),
	}
}

// Debug logs a debug message with structured fields
func (l *ZapLogger) Debug(ctx context.Context, message string, fields Fields) {
	l.logWithLevel(ctx, zapcore.DebugLevel, message, nil, fields)
}

// Info logs an info message with structured fields
func (l *ZapLogger) Info(ctx context.Context, message string, fields Fields) {
	l.logWithLevel(ctx, zapcore.InfoLevel, message, nil, fields)
}

// Warn logs a warning message with structured fields
func (l *ZapLogger) Warn(ctx context.Context, message string, fields Fields) {
	l.logWithLevel(ctx, zapcore.WarnLevel, message, nil, fields)
}

// Error logs an error message with error details and structured fields
func (l *ZapLogger) Error(ctx context.Context, message string, err error, fields Fields) {
	if fields == nil {
		fields = make(Fields)
	}
	if err != nil {
		fields["error"] = sanitizeLogString(err.Error())
	}
	l.logWithLevel(ctx, zapcore.ErrorLevel, message, err, fields)
}

// Fatal logs a fatal message and exits the program
func (l *ZapLogger) Fatal(ctx context.Context, message string, err error, fields Fields) {
	if fields == nil {
		fields = make(Fields)
	}
	if err != nil {
		fields["error"] = sanitizeLogString(err.Error())
	}
	l.logWithLevel(ctx, zapcore.FatalLevel, message, err, fields)
	os.Exit(1)
}

// WithContext adds context values to the logger
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	newLogger := &ZapLogger{
		logger: l.logger,
		fields: make(Fields),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add context values
	if ctx != nil {
		if requestID := getContextValue(ctx, "request_id"); requestID != "" {
			newLogger.fields["request_id"] = requestID
		}
		if userID := getContextValue(ctx, "user_id"); userID != "" {
			newLogger.fields["user_id"] = userID
		}
		if operation := getContextValue(ctx, "operation"); operation != "" {
			newLogger.fields["operation"] = operation
		}
	}

	return newLogger
}

// WithFields adds additional fields to the logger
func (l *ZapLogger) WithFields(fields Fields) Logger {
	newLogger := &ZapLogger{
		logger: l.logger,
		fields: make(Fields),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields (overwrite if key exists)
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// logWithLevel is the internal method that handles the actual logging
func (l *ZapLogger) logWithLevel(ctx context.Context, level zapcore.Level, message string, err error, fields Fields) {
	// Create Zap fields slice
	zapFields := make([]zap.Field, 0, len(fields)+2)

	// Add message field
	zapFields = append(zapFields, zap.String("message", sanitizeLogString(message)))

	// Add level field
	zapFields = append(zapFields, zap.String("level", level.String()))

	// Add all fields
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			zapFields = append(zapFields, zap.String(k, sanitizeLogString(val)))
		case int:
			zapFields = append(zapFields, zap.Int(k, val))
		case int64:
			zapFields = append(zapFields, zap.Int64(k, val))
		case float64:
			zapFields = append(zapFields, zap.Float64(k, val))
		case bool:
			zapFields = append(zapFields, zap.Bool(k, val))
		case error:
			zapFields = append(zapFields, zap.Error(val))
		default:
			zapFields = append(zapFields, zap.Any(k, val))
		}
	}

	// Add context fields if available
	if ctx != nil {
		if requestID := getContextValue(ctx, "request_id"); requestID != "" {
			zapFields = append(zapFields, zap.String("request_id", requestID))
		}
		if userID := getContextValue(ctx, "user_id"); userID != "" {
			zapFields = append(zapFields, zap.String("user_id", userID))
		}
		if operation := getContextValue(ctx, "operation"); operation != "" {
			zapFields = append(zapFields, zap.String("operation", operation))
		}
	}

	// Log with appropriate level
	switch level {
	case zapcore.DebugLevel:
		l.logger.Debug(message, zapFields...)
	case zapcore.InfoLevel:
		l.logger.Info(message, zapFields...)
	case zapcore.WarnLevel:
		l.logger.Warn(message, zapFields...)
	case zapcore.ErrorLevel:
		l.logger.Error(message, zapFields...)
	case zapcore.FatalLevel:
		l.logger.Fatal(message, zapFields...)
	}
}

// sanitizeLogString removes newline and carriage return characters to prevent log injection
func sanitizeLogString(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", "")
}

// getContextValue safely extracts a value from context
func getContextValue(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	if val := ctx.Value(key); val != nil {
		if str, ok := val.(string); ok {
			return sanitizeLogString(str)
		}
	}
	return ""
}

// Global logger instance
var globalLogger Logger

// Init initializes the global logger
func Init(level Level) {
	globalLogger = New(level)
}

// InitProduction initializes the global logger with production settings
func InitProduction() {
	globalLogger = NewProduction()
}

// InitDevelopment initializes the global logger with development settings
func InitDevelopment() {
	globalLogger = NewDevelopment()
}

// Get returns the global logger instance
func Get() Logger {
	if globalLogger == nil {
		globalLogger = New(INFO)
	}
	return globalLogger
}

// Convenience functions for global logger
func Debug(ctx context.Context, message string, fields Fields) {
	Get().Debug(ctx, message, fields)
}

func Info(ctx context.Context, message string, fields Fields) {
	Get().Info(ctx, message, fields)
}

func Warn(ctx context.Context, message string, fields Fields) {
	Get().Warn(ctx, message, fields)
}

func Error(ctx context.Context, message string, err error, fields Fields) {
	Get().Error(ctx, message, err, fields)
}

func Fatal(ctx context.Context, message string, err error, fields Fields) {
	Get().Fatal(ctx, message, err, fields)
}
