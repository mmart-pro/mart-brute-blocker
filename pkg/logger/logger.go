package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger(level, logfile string) (*Logger, error) {
	config := zap.NewProductionConfig()
	if logfile != "" {
		config.OutputPaths = append(config.OutputPaths, logfile)
	}
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(time.RFC3339))
	}

	zlevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("unknown level %s: %w", level, err)
	}
	config.Level = zlevel

	zlogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("can't build logger: %w", err)
	}

	return &Logger{logger: zlogger.Sugar()}, nil
}

func (l Logger) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l Logger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l Logger) Errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l Logger) Fatalf(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

func (l Logger) Close() {
	l.logger.Sync()
}
