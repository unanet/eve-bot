package tracer

import (
	"fmt"
	"io"

	"gitlab.unanet.io/devops/eve-bot/internal/evelogger"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

// Provider providers the tracer components
type Provider struct {
	Tracer opentracing.Tracer
	Closer io.Closer
}

type jaegerLoggerAdapter struct {
	logger evelogger.Container
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Bg().Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Bg().Info(fmt.Sprintf(msg, args...))
}

// New initializes the jaeger tracer
func New(name string, logger evelogger.Container, isGlobal bool) Provider {
	logger.Bg().With(zap.String("package", "tracer"))
	jCfg, err := jaegercfg.FromEnv()
	if err != nil {
		logger.Bg().Fatal("failed to parse Jaeger env vars", zap.Error(err))
	}

	jCfg.ServiceName = name
	jCfg.Sampler.Type = "const"
	jCfg.Sampler.Param = 1

	jaegerLogger := jaegerLoggerAdapter{logger: logger}
	// metricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: cfg.ServiceName, Tags: nil})

	tracer, closer, err := jCfg.NewTracer(
		jaegercfg.Logger(jaegerLogger),
		// jaegercfg.Metrics(metricsFactory),
	)

	if err != nil {
		logger.Bg().Fatal("cannot init Jaeger:", zap.Error(err))
	}

	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	if isGlobal {
		opentracing.SetGlobalTracer(tracer)
	}

	return Provider{Tracer: tracer, Closer: closer}

}
