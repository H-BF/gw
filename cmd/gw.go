package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	api "github.com/H-BF/gw/internal/api/SecGroup"
	"github.com/H-BF/gw/internal/authprovider"
	"github.com/H-BF/gw/internal/httpserver"
	"github.com/H-BF/gw/internal/telemetry/metrics"
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// port - the port var must be stored in config
const port = 8080

func main() {
	ctx := context.Background()

	casbinAuthProvider, err := authprovider.NewCasbinAuthProvider("model.conf", "policy.csv")
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	var sgOpts []connect.HandlerOption

	metrics.WhenMetricsEnabled(func(gm *metrics.GwMetrics) {
		metricsInterceptor, err := metrics.NewMetricInterceptor()
		if err != nil {
			log.Fatalln(err)
		}

		sgOpts = append(sgOpts,
			connect.WithInterceptors(metricsInterceptor),
			connect.WithRecover(func(ctx context.Context, spec connect.Spec, header http.Header, recInfo any) error {
				var service, method, clientName string // todo: заполнить переменные
				gm.ObservePanic(service, method, clientName)

				return nil
			}))
	})

	err = metrics.SetupMetric(context.Background(), func(reg *prometheus.Registry) error {
		mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	mux.Handle(sgroupsconnect.NewSecGroupServiceHandler(
		api.NewSecGroupService(casbinAuthProvider),
		sgOpts...,
	))

	if err = httpserver.ListenAndServe(port, mux); err != nil {
		log.Fatalln(err)
	}

	defer func(ctx context.Context) {
		if err := httpserver.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}(ctx)
}
