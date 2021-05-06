package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLog(level zapcore.Level) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.StacktraceKey = ""
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zap.Config{
		Encoding:      "json",
		Level:         zap.NewAtomicLevelAt(level),
		OutputPaths:   []string{"stdout"},
		EncoderConfig: config,
	}.Build()

	if err != nil {
		return nil, err
	}

	return logger, nil
}
