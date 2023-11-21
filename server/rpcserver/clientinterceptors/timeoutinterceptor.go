package clientinterceptors

/**
* created by mengqi on 2023/11/21
 */

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
