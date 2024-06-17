package grpcmetrics

import (
	"github.com/H-BF/gw/internal/telemetry/metrics/options"
	"github.com/prometheus/client_golang/prometheus"
)

type ConnMetrics struct {
	connGauge *prometheus.GaugeVec
}

func NewConnMetrics(serverOpts options.ServerMetricsOptions) *ConnMetrics {
	labels := [...]string{LabelRemoteAddr}
	opts := prometheus.GaugeOpts{
		Namespace: serverOpts.Namespace,
		Subsystem: serverOpts.Subsystem,
		Name:      "connections",
		Help:      "connection count at moment on a server",
	}

	return &ConnMetrics{
		connGauge: prometheus.NewGaugeVec(opts, labels[:]),
	}
}

func (cm *ConnMetrics) IncConn(remoteAddr string) {
	cm.connGauge.WithLabelValues(remoteAddr).Inc()
}

func (cm *ConnMetrics) DecConn(remoteAddr string) {
	cm.connGauge.WithLabelValues(remoteAddr).Inc()
}

func (cm *ConnMetrics) GetConnGauge() *prometheus.GaugeVec {
	return cm.connGauge
}
