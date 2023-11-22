package httpserver

/**
* created by mengqi on 2023/11/21
 */

type ServerOption func(*Server)

func WithEnableProfiling(profiling bool) ServerOption {
	return func(s *Server) {
		s.enableProfiling = profiling
	}
}

func WithMode(mode string) ServerOption {
	return func(s *Server) {
		s.mode = mode
	}
}

func WithServiceName(srvName string) ServerOption {
	return func(s *Server) {
		s.serviceName = srvName
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithMiddlewares(middlewares []string) ServerOption {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

func WithHealthz(healthz bool) ServerOption {
	return func(s *Server) {
		s.healthz = healthz
	}
}

func WithJwt(jwt *JwtInfo) ServerOption {
	return func(s *Server) {
		s.jwt = jwt
	}
}

func WithTransNames(transName string) ServerOption {
	return func(s *Server) {
		s.transName = transName
	}
}

func WithMetrics(enable bool) ServerOption {
	return func(o *Server) {
		o.enableMetrics = enable
	}
}
