package metrics

import (
	"context"

	"github.com/H-BF/gw/internal/telemetry/metrics/grpcmetrics"
	"github.com/H-BF/gw/internal/telemetry/metrics/options"
	"github.com/prometheus/client_golang/prometheus"
)

const ( // config constants
	// MetricEnable - time const-blank for metrics activation
	MetricEnable = true
)

type GwMetrics struct {
	*grpcmetrics.ConnMetrics
	*grpcmetrics.TotalRequestsMetric
	*grpcmetrics.ResponseTimeMetrics
	*grpcmetrics.GrpcErrorMetrics
}

/*
TODO: не нашел возможности включить/отключить сбор метрик го-рантайма в sgroups это делается с помощью опции `NoStandardMetrics`
*/

var (
	gmMetrics             *GwMetrics
	WhenMetricsEnabledFns []func(metrics *GwMetrics)
)

func SetupMetric(_ context.Context, f func(reg *prometheus.Registry) error, opts ...options.Options) error {
	if !MetricEnable {
		return nil
	}

	defaultServerMetricOptions := options.DefaultServerMetricsOptions()
	for _, opt := range opts {
		opt(&defaultServerMetricOptions)
	}

	gmMetrics = newGmMetrics(defaultServerMetricOptions)

	registry := prometheus.NewRegistry()
	collectors := []prometheus.Collector{
		gmMetrics.GetGrpcErrCounter(),
		gmMetrics.GetConnGauge(),
		gmMetrics.GetResTimeHist(),
	}

	reqTotalCollectors := gmMetrics.GetAllTotalRequestCollectors()
	collectors = append(collectors, reqTotalCollectors...)

	if len(defaultServerMetricOptions.StandardMetrics) > 0 {
		collectors = append(collectors, defaultServerMetricOptions.StandardMetrics...)
	}

	for _, collector := range collectors {
		if err := registry.Register(collector); err != nil {
			return err
		}
	}

	for _, f := range WhenMetricsEnabledFns {
		f(gmMetrics)
	}

	return f(registry)
}

func newGmMetrics(opts options.ServerMetricsOptions) *GwMetrics {
	return &GwMetrics{
		GrpcErrorMetrics:    grpcmetrics.NewGrpcErrorMetrics(opts),
		ConnMetrics:         grpcmetrics.NewConnMetrics(opts),
		TotalRequestsMetric: grpcmetrics.NewTotalRequestsMetric(opts),
		ResponseTimeMetrics: grpcmetrics.NewResponseTimeMetrics(opts),
	}
}

func GetGmMEtrics() *GwMetrics {
	return gmMetrics
}

func WhenMetricsEnabled(f func(gm *GwMetrics)) {
	if gmMetrics != nil {
		panic("should be called before `SetupMetric`")
	}
	WhenMetricsEnabledFns = append(WhenMetricsEnabledFns, f)
}
