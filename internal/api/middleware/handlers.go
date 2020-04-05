package middleware

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.unanet.io/devops/eve-bot/internal/evelogger"
	"gitlab.unanet.io/devops/eve-bot/internal/metrics"
	"gitlab.unanet.io/devops/eve-bot/internal/servicefactory"

	"github.com/opentracing-contrib/go-gorilla/gorilla"
	nethttp "github.com/opentracing-contrib/go-stdlib/nethttp"
	"go.uber.org/zap"
)

// TimeoutHandler will timeout the request if the TokenSvcTimeout is exceeded (currently 30 seconds)
// This will result in a 503 HTTP Error
func TimeoutHandler(timeoutSecs uint16) Adapter {
	return func(h http.Handler) http.Handler {
		return http.TimeoutHandler(h, time.Duration(timeoutSecs)*time.Second, fmt.Sprintf("Sorry! Service Timeout Exceeded: %v", timeoutSecs))
	}
}

// TracingHandler will adding jaeger tracing to all http requests
func TracingHandler(appCtxProvider *servicefactory.Container) Adapter {
	return func(h http.Handler) http.Handler {
		return gorilla.Middleware(
			appCtxProvider.TraceProvider.Tracer,
			h,
			nethttp.OperationNameFunc(func(r *http.Request) string {
				return "HTTP " + r.Method + " " + r.RequestURI
			}),
		)
	}
}

// LogMetricsHandler adapts the incoming request with Logging/Metrics
func LogMetricsHandler(logger evelogger.Container, metricProvider *metrics.Provider) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Establish new Context with Request-ID, StartTime, UserIP, etc.
			// continue to add request scoped "contextual data" here
			ctx := newWrappedReqCtx(r)

			// Log the incoming Request
			logger.For(ctx).Info("incoming request",
				zap.Time("request_start", RequestStartTimeFromContext(ctx)),
				zap.String("request_id", RequestIDFromContext(ctx)),
				zap.String("user_agent", r.UserAgent()),
				zap.Int64("content_length_bytes", r.ContentLength),
				zap.String("host", r.Host),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("user_ip", UserIPFromContext(ctx).String()),
				zap.String("protocol", r.Proto),
			)
			// Tally the incoming request metrics
			metricProvider.StatHTTPRequestCount.WithLabelValues(r.RequestURI, r.Method, r.Proto).Inc()
			metricProvider.StatRequestSaturationGuage.WithLabelValues(r.RequestURI, r.Method, r.Proto).Inc()

			// Run this on the way out (i.e. outgoing response)
			defer func() {
				// Calculate the request duration (i.e. latency)
				ms := float64(time.Since(RequestStartTimeFromContext(ctx))) / float64(time.Millisecond)

				// Log the outgoing response
				logger.For(ctx).Info("outgoing response",
					zap.String("request_id", RequestIDFromContext(ctx)),
					zap.Float64("duration_ms", ms),
					zap.String("user_agent", r.UserAgent()),
					zap.Int64("content_length_bytes", r.ContentLength),
					zap.String("host", r.Host),
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("uri", r.RequestURI),
					zap.String("method", r.Method),
					zap.String("user_ip", UserIPFromContext(ctx).String()),
					zap.String("protocol", r.Proto),
				)

				// Tally the outgoing response metrics
				metricProvider.StatRequestDurationHistogram.WithLabelValues(r.RequestURI, r.Method, r.Proto).Observe(ms)
				metricProvider.StatRequestDurationGuage.WithLabelValues(r.RequestURI, r.Method, r.Proto).Set(ms)
				metricProvider.StatRequestSaturationGuage.WithLabelValues(r.RequestURI, r.Method, r.Proto).Dec()

			}()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RouteHandlerSig is the route handler signature
type RouteHandlerSig func(appCtxProvider *servicefactory.Container, res http.ResponseWriter, req *http.Request) (int, interface{}, error)

// Handler is the wrapper that provides context to the app handler
type Handler struct {
	AppCtx       *servicefactory.Container
	RouteHandler RouteHandlerSig
}
