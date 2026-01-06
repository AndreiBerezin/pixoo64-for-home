package log

import (
	"context"

	"go.uber.org/zap"
)

var defaultLogger Logger

func DefaultLogger() Logger {
	return defaultLogger
}

func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}

func With(fields ...zap.Field) Logger {
	return defaultLogger.With(fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger.l.Info(msg, append(fields, defaultLogger.fields...)...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger.l.Warn(msg, append(fields, defaultLogger.fields...)...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger.l.Error(msg, append(fields, defaultLogger.fields...)...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.l.Fatal(msg, append(fields, defaultLogger.fields...)...)
}

func Debug(msg string, fields ...zap.Field) {
	defaultLogger.l.Debug(msg, append(fields, defaultLogger.fields...)...)
}

func Sync() {
	_ = defaultLogger.l.Sync()
}
