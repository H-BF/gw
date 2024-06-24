package grpcmetrics

import (
	"github.com/H-BF/gw/internal/telemetry/metrics/options"
	"github.com/prometheus/client_golang/prometheus"
)

type TotalRequestsMetric struct {
	messages       *prometheus.CounterVec
	methodStarted  *prometheus.CounterVec
	methodFinished *prometheus.CounterVec
	methodPanicked *prometheus.CounterVec
}

const (
	Received = "received"
	Sent     = "sent"
)

func NewTotalRequestsMetric(serverMetricsOptions options.ServerMetricsOptions) *TotalRequestsMetric {
	messages := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverMetricsOptions.Namespace,
		Subsystem: serverMetricsOptions.Subsystem,
		Name:      "messages",
		Help:      "received and sent message counters",
	}, []string{LabelService, LabelMethod, LabelState})

	started := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverMetricsOptions.Namespace,
		Subsystem: serverMetricsOptions.Subsystem,
		Name:      "methods_started",
		Help:      "started methods counter",
	}, []string{LabelService, LabelMethod, LabelClientName})

	finished := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverMetricsOptions.Namespace,
		Subsystem: serverMetricsOptions.Subsystem,
		Name:      "methods_finished",
		Help:      "finished methods counter",
	}, []string{LabelService, LabelMethod, LabelClientName, LabelGRPCCode})

	panicked := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverMetricsOptions.Namespace,
		Subsystem: serverMetricsOptions.Subsystem,
		Name:      "methods_panicked",
		Help:      "panicked methods counter",
	}, []string{LabelService, LabelMethod, LabelClientName})

	return &TotalRequestsMetric{
		messages:       messages,
		methodStarted:  started,
		methodFinished: finished,
		methodPanicked: panicked,
	}

}

func (trm *TotalRequestsMetric) GetAllTotalRequestCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		trm.messages,
		trm.methodFinished,
		trm.methodStarted,
		trm.methodPanicked,
	}
}

func (trm *TotalRequestsMetric) IncStartedMethod(service, method, clientName string) {
	labels := prometheus.Labels{
		LabelService:    service,
		LabelMethod:     method,
		LabelClientName: clientName,
	}

	trm.methodStarted.With(labels).Inc()
}

func (trm *TotalRequestsMetric) IncFinishedMethod(service, method, clientName, gRPCCode string) {
	labels := prometheus.Labels{
		LabelService:    service,
		LabelMethod:     method,
		LabelClientName: clientName,
		LabelGRPCCode:   gRPCCode,
	}

	trm.methodFinished.With(labels).Inc()
}

func (trm *TotalRequestsMetric) IncReceivedSentMessage(service, method string, isClient bool) {
	requestSpan := Sent

	// todo: клиент может только отправлять сообщения, а сервер - принимать??? подумать над этим и исправить логику
	if !isClient {
		requestSpan = Received
	}

	labels := prometheus.Labels{
		LabelService: service,
		LabelMethod:  method,
		LabelState:   requestSpan,
	}

	trm.messages.With(labels).Inc()
}

func (trm *TotalRequestsMetric) ObservePanic(service, method, clientName string) {
	labels := prometheus.Labels{
		LabelService:    service,
		LabelMethod:     method,
		LabelClientName: clientName,
	}

	trm.methodPanicked.With(labels).Inc()
}
