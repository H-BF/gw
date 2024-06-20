package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/grpc/status"
)

func NewMetricInterceptor() (connect.Interceptor, error) {
	var metricsInterceptor connect.UnaryInterceptorFunc = func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if !MetricEnable {
				return next(ctx, req)
			}

			gm := GetGmMEtrics()

			gm.IncConn(req.Peer().Addr)
			defer gm.DecConn(req.Peer().Addr)

			procedure := strings.Split(req.Spec().Procedure, "/")
			service := procedure[1]
			method := procedure[2]

			clientName := req.Header().Get("user-agent")

			// todo: in the future put it'n a sep interface `panicObserver`
			defer func(gm *GwMetrics, service, method, clientName string) {
				if panicked := recover(); panicked != nil {
					gm.ObservePanic(service, method, clientName)
					fmt.Println(panicked)
				}
			}(gm, service, method, clientName)

			gm.IncReceivedSentMessage(service, method, req.Spec().IsClient)

			start := time.Now()
			res, err := next(ctx, req)
			if err != nil {
				gm.IncGrpcErrCount(service, method, err.Error())
			}

			gm.ObserveResTime(service, method, float64(time.Since(start).Microseconds()))

			grpcCode := status.Code(err).String()

			gm.IncStartedMethod(service, method, clientName)
			defer gm.IncFinishedMethod(service, method, clientName, grpcCode)

			return res, err
		})
	}

	return connect.UnaryInterceptorFunc(metricsInterceptor), nil
}
