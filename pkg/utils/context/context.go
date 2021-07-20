package context

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

func NewLoggingContext(level zapcore.Level) context.Context {
	return logging.WithLogger(signals.NewContext(), loggerAt(level))
}

func loggerAt(level zapcore.Level) *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.DisableCaller = true
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
