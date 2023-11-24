package serverinterceptors

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/chaos-ma/chaos/core/metric"
)

const serverNamespace = "rpc_server"

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

func UnaryPrometheusInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	startTime := time.Now()
	resp, err = handler(ctx, req)

	//记录了耗时
	metricServerReqDur.Observe(int64(time.Since(startTime)/time.Millisecond), info.FullMethod)

	//记录了状态码
	metricServerReqCodeTotal.Inc(info.FullMethod, strconv.Itoa(int(status.Code(err))))
	return resp, err
}
