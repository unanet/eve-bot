package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
)

// Provider the metrics provider which holds all metrics
type Provider struct {
	StatMemAllocGuage, StatMemTotalAllocGuage, StatMemSysGuage, StatMemNumGCGuage, StatGoRoutineGuage prometheus.Gauge
	StatRequestSaturationGuage, StatRequestDurationGuage, StatBuildInfo                               *prometheus.GaugeVec
	StatHTTPRequestCount, StatHTTPResponseCount, StatAuditCount                                       *prometheus.CounterVec
	StatRequestDurationHistogram                                                                      *prometheus.HistogramVec
}

// New returns an initialized Metric Provider
func New() *Provider {
	return &Provider{
		StatMemAllocGuage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_allocation",
				Help: "The memory allocation guage (runtime.MemStats.Alloc)",
			}),
		StatMemTotalAllocGuage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_total_alloc",
				Help: "The total memory allocation guage (runtime.MemStats.TotalAlloc)",
			}),
		StatMemSysGuage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_system",
				Help: "The sys memory guage (runtime.MemStats.Sys)",
			}),
		StatMemNumGCGuage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gc_count",
				Help: "The garbage colletion gauge (runtime.MemStats.NumGC)",
			}),
		StatGoRoutineGuage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "go_routine_guage",
				Help: "The number of go routines in the runtime (runtime.NumGoroutine)",
			}),
		StatRequestSaturationGuage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_request_saturation",
				Help: "The total number of requests inside the server (transactions serving)",
			}, []string{"uri", "method", "protocol"}),
		StatBuildInfo: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "service_build_info",
				Help: "A metric with a constant '1' value labeled by version, revision, branch, and goversion from which the service was build was built.",
			}, []string{"service", "revision", "branch", "version", "author", "build_date", "build_user", "build_host"}),
		StatRequestDurationGuage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_request_duration_ms",
				Help: "The time the server spends processing a request in milliseconds",
			}, []string{"uri", "method", "protocol"}),
		StatAuditCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "audit_total",
				Help: "The total number of audit events",
			}, []string{"event"}),
		StatHTTPRequestCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "The total number of incoming requests to the service",
			}, []string{"uri", "method", "protocol"}),
		StatHTTPResponseCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_response_total",
				Help: "The total number of outgoing responses to the client",
			}, []string{"code", "uri", "method", "protocol"}),
		StatRequestDurationHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_histogram_ms",
				Help:    "time spent processing an http request in milliseconds",
				Buckets: prometheus.ExponentialBuckets(0.1, 2, 18),
			}, []string{"uri", "method", "protocol"}),
	}
}

// LogLevelPrometheusHook created the PrometheusHook for log level counts
func LogLevelPrometheusHook() *PrometheusHook {
	return &PrometheusHook{
		counter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "log_level_total",
				Help: "Number of log statements, differentiated by log level.",
			},
			[]string{"level"},
		),
	}
}

// PrometheusHook contains the counter
type PrometheusHook struct {
	counter *prometheus.CounterVec
}

// Levels returns all of the log levels
func (h *PrometheusHook) Levels() []zapcore.Level {
	return []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.PanicLevel}
}

// Fire runs the prometheus counter
func (h *PrometheusHook) Fire(e zapcore.Entry) error {
	h.counter.WithLabelValues(e.Level.String()).Inc()
	return nil
}
