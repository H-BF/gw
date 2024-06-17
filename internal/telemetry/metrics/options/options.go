package options

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type (
	ServerMetricsOptions struct {
		Namespace string
		Subsystem string

		StandardMetrics []prometheus.Collector
	}

	Options func(*ServerMetricsOptions)
)

const (
	DefaultNamespace = "sys"
	DefaultSubsystem = "grpc_server"
)

func DefaultServerMetricsOptions() ServerMetricsOptions {
	return ServerMetricsOptions{
		Namespace: DefaultNamespace,
		Subsystem: DefaultSubsystem,
	}
}

func WithNamespace(ns string) Options {
	return func(smo *ServerMetricsOptions) {
		smo.Namespace = ns
	}
}

func WithSubsystem(ss string) Options {
	return func(smo *ServerMetricsOptions) {
		smo.Subsystem = ss
	}
}

func WithStandardMetrics() Options {
	return func(smo *ServerMetricsOptions) {
		smo.StandardMetrics = []prometheus.Collector{
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		}
	}
}
