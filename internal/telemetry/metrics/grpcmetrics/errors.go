package grpcmetrics

import (
	"github.com/H-BF/gw/internal/telemetry/metrics/options"
	"github.com/prometheus/client_golang/prometheus"
)

type GrpcErrorMetrics struct {
	allErrorsCounterVec *prometheus.CounterVec
}

func NewGrpcErrorMetrics(serverOpts options.ServerMetricsOptions) *GrpcErrorMetrics {
	labels := [...]string{LabelService, LabelMethod, LabelErrSrc}
	opts := prometheus.CounterOpts{
		Namespace: serverOpts.Namespace,
		Subsystem: serverOpts.Subsystem,
		Name:      "errors",
		Help:      "count of errors",
	}

	return &GrpcErrorMetrics{
		allErrorsCounterVec: prometheus.NewCounterVec(opts, labels[:]),
	}
}

func (gem *GrpcErrorMetrics) GetGrpcErrCounter() *prometheus.CounterVec {
	return gem.allErrorsCounterVec
}

func (gem *GrpcErrorMetrics) IncGrpcErrCount(service, method, err string) {
	labs := prometheus.Labels{
		LabelService: service,
		LabelMethod:  method,
		LabelErrSrc:  err,
	}

	gem.allErrorsCounterVec.With(labs).Inc()
}
