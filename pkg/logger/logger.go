package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	logger    *zap.Logger
	startTime time.Time
)

func ZapLoggerInit() {
	startTime = time.Now()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writer := zapcore.AddSync(os.Stderr)
	logLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writer, logLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func ServiceError(custom, err error, data ...any) {
	Error(fmt.Sprintf("%v, err - %v, info - %+v", custom, err, data))
}

func Uptime() time.Duration {
	return time.Since(startTime)
}
