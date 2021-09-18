package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Http_request_duration_seconds = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of lantencies for HTTP requests",
		},
	)

	Http_request_qps = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_request_qps",
			Help: "The number of HTTP requests on / served in the last second",
		},
	)

	Http_request_in_flight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_request_in_flight",
			Help: "Current number of http requests in flight",
		},
	)

	Http_request_total = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The total number of processed http requests",
		},
	)
)
