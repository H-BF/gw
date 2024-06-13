package metrics

import (
	"context"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

const ( // config constants
	// MetricEnable - time const-blank for metrics activation
	MetricEnable = true
	UserAgent    = "telemetry/user/gw"
)

type GwMetrics struct {
	errorCount   *prometheus.CounterVec
	grpcMessages *prometheus.CounterVec
}

const nsGw = "gw"

const (
	labelUserAgent   = "user_agent"
	labelHostName    = "host_name"
	labelSource      = "source"
	labelGrpcMethod  = "method"
	labelGrpcService = "service"
)

const ErrSrcGwServer = "gw-svc"

var gmMetrics *GwMetrics

func SetupMetric(ctx context.Context, f func(reg *prometheus.Registry) error) error {
	if !MetricEnable {
		return nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	labels := prometheus.Labels{
		labelUserAgent: UserAgent,
		labelHostName:  hostname,
	}

	gmMetrics = newGmMetrics(labels)

	registry := prometheus.NewRegistry()
	collectors := []prometheus.Collector{
		gmMetrics.errorCount,
		gmMetrics.grpcMessages,
	}

	for _, collector := range collectors {
		if err = registry.Register(collector); err != nil {
			return err
		}
	}

	return f(registry)
}

func newGmMetrics(labels prometheus.Labels) *GwMetrics {
	return &GwMetrics{
		errorCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   nsGw,
			Name:        "errors",
			Help:        "count of errors",
			ConstLabels: labels,
		}, []string{labelSource}),
		grpcMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   nsGw,
			Name:        "server_grpc_messages",
			Help:        "received and sent message counters",
			ConstLabels: labels,
		}, []string{labelGrpcMethod, labelGrpcService}),
	}
}

func GetGmMEtrics() *GwMetrics {
	return gmMetrics
}

func (gm *GwMetrics) IncError(src string) {
	gm.errorCount.WithLabelValues(src).Inc()
}

func (gm *GwMetrics) IncGrpcMessage(method, service string) {
	gm.grpcMessages.WithLabelValues(method, service).Inc()
}
