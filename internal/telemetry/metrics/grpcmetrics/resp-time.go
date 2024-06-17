package grpcmetrics

import (
	"github.com/H-BF/gw/internal/telemetry/metrics/options"
	"github.com/prometheus/client_golang/prometheus"
)

type ResponseTimeMetrics struct {
	respTimeHist *prometheus.HistogramVec
}

func NewResponseTimeMetrics(serverOpts options.ServerMetricsOptions) *ResponseTimeMetrics {
	rtm := &ResponseTimeMetrics{}

	labels := [...]string{LabelService, LabelMethod}
	opts := prometheus.HistogramOpts{
		Namespace: serverOpts.Namespace,
		Subsystem: serverOpts.Subsystem,
		Name:      "response_time",
		Help:      "response time duration in milliseconds",
		Buckets:   rtm.defaultBucket(),
	}

	rtm.respTimeHist = prometheus.NewHistogramVec(opts, labels[:][:])

	return rtm
}

func (*ResponseTimeMetrics) defaultBucket() []float64 {
	return []float64{
		.0001, .0005, .00075, .001, .0025, .005, 0.0075, .01, 0.025, .05, 0.075,
		.1, .25, .5, .75, 10, 25, 50, 75, 100, 500, 1000}
}

func (rtm *ResponseTimeMetrics) GetResTimeHist() *prometheus.HistogramVec {
	return rtm.respTimeHist
}

func (rtm *ResponseTimeMetrics) ObserveResTime(service, method string, ms float64) {
	labs := prometheus.Labels{
		LabelMethod:  method,
		LabelService: service,
	}

	rtm.respTimeHist.With(labs).Observe(ms)
}
