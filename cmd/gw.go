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

	metricsInterceptor, err := metrics.NewMetricInterceptor()
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.Handle(sgroupsconnect.NewSecGroupServiceHandler(
		api.NewSecGroupService(casbinAuthProvider),
		connect.WithInterceptors(metricsInterceptor),
	))

	if err := metrics.SetupMetric(context.Background(), func(reg *prometheus.Registry) error {
		mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
		return nil
	}); err != nil {
		log.Fatalln(err)
	}

	if err = httpserver.ListenAndServe(port, mux); err != nil {
		log.Fatalln(err)
	}

	defer func(ctx context.Context) {
		if err := httpserver.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}(ctx)
}

func init() {

}
