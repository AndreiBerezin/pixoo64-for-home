package log

import (
	"context"
	"os"

	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	l      *zap.Logger
	fields []zap.Field
}

func Init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "@timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "line_number"
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writer := zapcore.Lock(os.Stdout)

	level := zapcore.InfoLevel
	if env.IsDebug() {
		level = zapcore.DebugLevel
	}
	core := zapcore.NewCore(encoder, writer, level)

	zapLogger := zap.New(core)
	zapLogger = zapLogger.With(zap.String("env", os.Getenv("ENV")))

	defaultLogger = Logger{l: zapLogger, fields: []zap.Field{}}
}

func (l Logger) WithContext(ctx context.Context) Logger {
	fields := ctx.Value(ctxFieldsKey{})
	if fields == nil {
		return l
	}
	newl := Logger{l: l.l, fields: append(append([]zap.Field{}, fields.([]zap.Field)...), l.fields...)}
	return newl
}

func (l Logger) With(fields ...zap.Field) Logger {
	newl := Logger{l: l.l, fields: append(fields, l.fields...)}
	return newl
}

func (l Logger) Info(msg string, fields ...zap.Field) {
	l.l.Info(msg, append(fields, l.fields...)...)
}

func (l Logger) Warn(msg string, fields ...zap.Field) {
	l.l.Warn(msg, append(fields, l.fields...)...)
}

func (l Logger) Error(msg string, fields ...zap.Field) {
	l.l.Error(msg, append(fields, l.fields...)...)
}

func (l Logger) Fatal(msg string, fields ...zap.Field) {
	l.l.Fatal(msg, append(fields, l.fields...)...)
}

func (l Logger) Debug(msg string, fields ...zap.Field) {
	l.l.Debug(msg, append(fields, l.fields...)...)
}
