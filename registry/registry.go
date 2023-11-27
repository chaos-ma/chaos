package registry

/**
* created by mengqi on 2023/11/13
 */

import (
	"context"
)

type Registrar interface {
	Register(ctx context.Context, service *ServiceInstance) error
	Deregister(ctx context.Context, service *ServiceInstance) error
}

type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	Next() ([]*ServiceInstance, error) //1. 第一次监听时，如果服务实例列表不为空，则返回服务实例列表。2. 如果服务实例发生变化，则返回服务实例列表。3. 如果上面两种情况都不满足，则会阻塞到context deadline或者cancel
	Stop() error
}

type ServiceInstance struct {
	ID        string            `json:"id"`        //注册到注册中心的服务id
	Name      string            `json:"name"`      //服务名称
	Version   string            `json:"version"`   //服务版本
	Metadata  map[string]string `json:"metadata"`  //服务元数据
	Endpoints []string          `json:"endpoints"` //http://127.0.0.1:8000 ,grpc://127.0.0.1:9000
}
