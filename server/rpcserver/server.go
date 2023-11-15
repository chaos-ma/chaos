package rpcserver

/**
* created by mengqi on 2023/11/15
 */

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"net"
	"net/url"
)

type ServerOption func(o *RpcServer)

type RpcServer struct {
	*grpc.Server
	address            string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	grpcOpts           []grpc.ServerOption
	listener           net.Listener
	health             *health.Server
	endpoint           *url.URL
}

func WithAddress(address string) ServerOption {
	return func(s *RpcServer) {
		s.address = address
	}
}

func WithListener(listener net.Listener) ServerOption {
	return func(s *RpcServer) {
		s.listener = listener
	}
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *RpcServer) {
		s.unaryInterceptors = interceptors
	}
}

func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *RpcServer) {
		s.streamInterceptors = interceptors
	}
}

func WithGrpcOpts(opts ...grpc.ServerOption) ServerOption {
	return func(s *RpcServer) {
		s.grpcOpts = opts
	}
}
