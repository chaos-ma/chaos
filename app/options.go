package app

/**
* created by mengqi on 2023/11/13
* 这里使用函数选项模式
 */

import (
	"net/url"
	"os"
	"time"

	"github.com/chaos-ma/chaos/registry"
	"github.com/chaos-ma/chaos/server/httpserver"
	"github.com/chaos-ma/chaos/server/rpcserver"
)

type Option func(o *options)

type options struct {
	id        string
	endpoints []*url.URL
	name      string

	sigs []os.Signal

	//允许用户传入自己的实现
	registrar        registry.Registrar
	registrarTimeout time.Duration

	//stop超时时间
	stopTimeout time.Duration

	restServer *httpserver.Server
	rpcServer  *rpcserver.Server
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithRPCServer(server *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = server
	}
}

func WithRestServer(server *httpserver.Server) Option {
	return func(o *options) {
		o.restServer = server
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
