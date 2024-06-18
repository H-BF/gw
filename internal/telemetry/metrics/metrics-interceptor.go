package metrics

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
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

			gm.IncReceivedSentMessage(service, method, req.Spec().IsClient)

			start := time.Now()
			res, err := next(ctx, req)
			if err != nil {
				gm.IncGrpcErrCount(service, method, err.Error())
			}

			gm.ObserveResTime(service, method, float64(time.Since(start).Microseconds()))

			clientName := req.Header().Get("user-agent")
			grpcCode := res.Trailer().Get("Grpc-Status")[0]

			gm.IncStartedMethod(service, method, clientName)
			defer gm.IncFinishedMethod(service, method, clientName, string(grpcCode))

			return res, err
		})
	}

	return connect.UnaryInterceptorFunc(metricsInterceptor), nil
}
