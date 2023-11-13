package app

/**
* created by mengqi on 2023/11/13
* 这里使用函数选项模式
 */

import (
	"chaos/registry"
	"net/url"
	"os"
	"time"
)

type Option func(o *serviceOptions)

type serviceOptions struct {
	id                string
	name              string
	endpoints         []*url.URL
	signals           []os.Signal              //监听的信号量
	registry          registry.ServiceRegistry //接受自定义的注册方法
	registerTimeout   time.Duration            //注册超时时间
	unRegisterTimeout time.Duration            //注销超时时间
}

func WithID(id string) Option {
	return func(o *serviceOptions) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *serviceOptions) {
		o.name = name
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *serviceOptions) {
		o.endpoints = endpoints
	}
}

func WithSignals(signals []os.Signal) Option {
	return func(o *serviceOptions) {
		o.signals = signals
	}
}

func WithRegistry(registry registry.ServiceRegistry) Option {
	return func(o *serviceOptions) {
		o.registry = registry
	}
}
