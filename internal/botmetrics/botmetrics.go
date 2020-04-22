package botmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// StatBotErrCount counter for Bot errors
	// We will want to fire an alarm any time these occur
	StatBotErrCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evebot_err_total",
			Help: "The total number of eve errors",
		}, []string{"err_type"})
)
