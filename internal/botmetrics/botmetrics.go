package botmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (

	// StatBotWorkReqReceivedCount tracks the count of incoming work requests
	// this is used for "async" calls to the API
	// Deployment/Migration where callback is being used
	StatBotWorkReqReceivedCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evebot_work_req_received_total",
			Help: "The total number of eve work requests",
		}, []string{"work_type"})

	StatBotWorkReqInProcessCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evebot_work_req_wip_total",
			Help: "The total number of eve work requests in process (WIP)",
		}, []string{"work_type", "worker_id"})

	StatBotWorkReqCompletedCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evebot_work_req_completed_total",
			Help: "The total number of completed eve work requests",
		}, []string{"work_type", "worker_id"})

	StatBotWorkerQueueSaturationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evebot_worker_queue_saturation",
			Help: "The number of work requests waiting to be processed by a worker",
		}, []string{"work_type"})

	StatBotWorkerWIPSaturationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evebot_worker_wip_saturation",
			Help: "The number of work requests currently being processed by the workers (WIP)",
		}, []string{"work_type", "worker_id"})

	StatBotWorkerQueueDurationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evebot_worker_queue_duration_ms",
			Help: "The work request spent waiting in queue to be processed milliseconds",
		}, []string{"work_type", "worker_id"})

	StatBotWorkerQueueDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "evebot_worker_queue_duration_histogram_ms",
			Help:    "time work request spent waiting in queue to be processed in milliseconds",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 18),
		}, []string{"work_type", "worker_id"})

	StatBotWorkerInProcessDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "evebot_worker_inprocess_duration_histogram_ms",
			Help:    "time spent processing work request in milliseconds",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 18),
		}, []string{"work_type", "worker_id"})

	StatBotWorkerInProcessDurationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evebot_worker_inprocess_duration_ms",
			Help: "The work request spent waiting in queue to be processed milliseconds",
		}, []string{"work_type", "worker_id"})
)
