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
	// Service    string
	// Method     string
	// ClientName string
	errorCount *prometheus.CounterVec
}

const (
	labelUserAgent = "user_agent"
	labelHostName  = "host_name"
	labelSource    = "source"
	nsGw           = "gw"
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
	}
}

func GetGmMEtrics() *GwMetrics {
	return gmMetrics
}

func (gm *GwMetrics) IncError(src string) {
	gm.errorCount.WithLabelValues(src).Inc()
}
