package evelogger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/opentracing/opentracing-go"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/metrics"
)

// Container conains the application logger
type Container struct {
	logger *zap.Logger
}

func logLevel(cfgLevel string) zap.AtomicLevel {
	var logLevel zap.AtomicLevel
	switch cfgLevel {
	case "debug":
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error", "err":
		logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		logLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return logLevel
}

func PrometheusHook(e zapcore.Entry) error {
	return nil
}

func NewLogContainer(botCfg *config.Config, logger *zap.Logger) Container {
	cfg := zap.Config{
		Level:            logLevel(botCfg.Logger.Level),
		Encoding:         botCfg.Logger.Encoding,
		OutputPaths:      botCfg.Logger.OutputPaths,
		ErrorOutputPaths: botCfg.Logger.ErrorOutputPaths,
		InitialFields: map[string]interface{}{
			"service": botCfg.API.ServiceName,
		},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zapcore.RegisterHooks(logger.Core(), metrics.LogLevelPrometheusHook().Fire)

	return Container{logger: logger}

}

// Bg creates a context-unaware logger.
func (b Container) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Container) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		// TODO for Jaeger span extract trace/span IDs as fields
		return spanLogger{span: span, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b Container) With(fields ...zapcore.Field) Container {
	return Container{logger: b.logger.With(fields...)}
}
