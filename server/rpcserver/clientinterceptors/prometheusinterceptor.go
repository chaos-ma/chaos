package clientinterceptors

/**
* created by mengqi on 2023/11/21
 */

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"chaos/core/metric"
)

const serverNamespace = "rpc_client"

/*
两个基本指标。 1. 每个请求的耗时(histogram) 2. 每个请求的状态计数器(counter)
/user 状态码 有label 主要是状态码
*/

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "chaos_duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "chaos_code_total",
		Help:      "rpc server requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func PrometheusInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		//记录了耗时
		metricServerReqDur.Observe(int64(time.Since(startTime)/time.Millisecond), method)

		//记录了状态码
		metricServerReqCodeTotal.Inc(method, strconv.Itoa(int(status.Code(err))))
		return err
	}
}
