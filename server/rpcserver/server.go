package rpcserver

/**
* created by mengqi on 2023/11/15
 */

import (
	"context"
	"net"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"chaos/log"
	apiMetaData "chaos/metadata"
	serverInterceptors "chaos/server/rpcserver/serverinterceptors"
	"chaos/utils/host"
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
	metadata           *apiMetaData.Server
	timeout            time.Duration
	enableMetrics      bool
}

func NewServer(opts ...ServerOption) *RpcServer {
	srv := &RpcServer{
		address: ":0",
		health:  health.NewServer(),
		//timeout: 1 * time.Second,
	}

	for _, o := range opts {
		o(srv)
	}
	// 默认拦截器crash，tracing
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverInterceptors.UnaryCrashInterceptor,
		otelgrpc.UnaryServerInterceptor(),
	}

	//加入监控拦截器
	if srv.enableMetrics {
		unaryInterceptors = append(unaryInterceptors, serverInterceptors.UnaryPrometheusInterceptor)
	}

	//加入超时拦截器
	if srv.timeout > 0 {
		unaryInterceptors = append(unaryInterceptors, serverInterceptors.UnaryTimeoutInterceptor(srv.timeout))
	}

	if len(srv.unaryInterceptors) > 0 {
		unaryInterceptors = append(unaryInterceptors, srv.unaryInterceptors...)
	}

	//传入的拦截器转换成grpc的ServerOption
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInterceptors...)}

	//把自定义的grpc.ServerOption放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	//注册metadata的Server
	srv.metadata = apiMetaData.NewServer(srv.Server)

	//解析address
	err := srv.listenAndEndpoint()
	if err != nil {
		panic(err)
	}

	//健康检查
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	apiMetaData.RegisterMetadataServer(srv.Server, srv.metadata)
	reflection.Register(srv.Server)
	//可以支持用户直接通过grpc的一个接口查看当前支持的所有的rpc服务

	return srv
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

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *RpcServer) {
		s.timeout = timeout
	}
}

// ip和端口的提取
func (s *RpcServer) listenAndEndpoint() error {
	if s.listener == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.listener = lis
	}
	addr, err := host.Extract(s.address, s.listener)
	if err != nil {
		_ = s.listener.Close()
		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

// Start 启动grpc的服务
func (s *RpcServer) Start(ctx context.Context) error {
	log.Infof("[grpc] server listening on: %s", s.listener.Addr().String())
	s.health.Resume()
	return s.Serve(s.listener)
}

func (s *RpcServer) Stop(ctx context.Context) error {
	//设置服务的状态为not_serving，防止接收新的请求过来
	s.health.Shutdown()
	s.GracefulStop()
	log.Infof("[grpc] server stopped")
	return nil
}
