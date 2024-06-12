package metrics

import (
	"context"

	"connectrpc.com/connect"
)

func NewMetricInterceptor() (connect.Interceptor, error) {
	var metricsInterceptor connect.UnaryInterceptorFunc = func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if !MetricEnable {
				return next(ctx, req)
			}

			gm := GetGmMEtrics()

			res, err := next(ctx, req)
			if err != nil {
				gm.IncError(ErrSrcGwServer)
			}

			return res, err
		})
	}

	return connect.UnaryInterceptorFunc(metricsInterceptor), nil
}
