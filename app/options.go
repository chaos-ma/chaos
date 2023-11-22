package app

/**
* created by mengqi on 2023/11/13
* 这里使用函数选项模式
 */

import (
	"net/url"
	"os"
	"time"

	"chaos/registry"
)

type Option func(o *options)

type options struct {
	id                string
	name              string
	endpoints         []*url.URL
	signals           []os.Signal       //监听的信号量
	registry          registry.Registry //接受自定义的注册方法
	registerTimeout   time.Duration     //注册超时时间
	unRegisterTimeout time.Duration     //注销超时时间
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

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithSignals(signals []os.Signal) Option {
	return func(o *options) {
		o.signals = signals
	}
}

func WithRegistry(registry registry.Registry) Option {
	return func(o *options) {
		o.registry = registry
	}
}
