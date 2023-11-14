package app

/**
* created by mengqi on 2023/11/13
 */

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"

	"chaos/registry"
)

type App struct {
	opts     serviceOptions
	instance *registry.ServiceInstance
	lock     sync.Mutex //保证获取instance的时候是线程安全
}

func NewApp(opts ...Option) *App {
	//设置默认值
	o := serviceOptions{
		signals:           []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registerTimeout:   10 * time.Second,
		unRegisterTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err != nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}
	return &App{opts: o}
}

func (app *App) Run() error {
	//创建实例
	instance, err := app.buildInstance()
	if err != nil {
		return err
	}
	app.lock.Lock()
	app.instance = instance
	app.lock.Unlock()
	//注册服务
	if app.opts.registry != nil {
		//设置超时时间
		regCtx, regCancel := context.WithTimeout(context.Background(), app.opts.registerTimeout)
		defer regCancel()
		err := app.opts.registry.Register(regCtx, instance)
		if err != nil {
			// TODO 打印日志信息
			return err
		}
	}

	//监听退出信号量
	c := make(chan os.Signal, 1)
	signal.Notify(c, app.opts.signals...)
	<-c
	return nil
}

func (app *App) Stop() error {
	app.lock.Lock()
	instance := app.instance
	app.lock.Unlock()
	//注销服务
	if app.opts.registry != nil && instance != nil {
		//设置超时时间
		regCtx, regCancel := context.WithTimeout(context.Background(), app.opts.unRegisterTimeout)
		defer regCancel()
		err := app.opts.registry.UnRegister(regCtx, instance)
		if err != nil {
			// TODO 打印日志信息
			return err
		}
	}
	return nil
}

func (app *App) Restart() error {
	return nil
}

func (app *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range app.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	return &registry.ServiceInstance{
		ID:        app.opts.id,
		Name:      app.opts.name,
		Endpoints: endpoints,
	}, nil
}
