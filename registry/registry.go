package registry

/**
* created by mengqi on 2023/11/13
 */

import "context"

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	Register(ctx context.Context, service *ServiceInstance) error   //注册服务
	UnRegister(ctx context.Context, service *ServiceInstance) error //注销服务
}

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error) //获取服务实例
	Watch(ctx context.Context, serviceName string) (ServiceWatcher, error)          //
}

// ServiceWatcher 服务监听接口
type ServiceWatcher interface {
	Next([]*ServiceInstance, error) //获取服务实例 1.服务实例不为空，返回实例。2.服务实例发生变化，返回实例。3.如果都不满足，则阻塞context deadline
	Stop() error                    //放弃监听
}

type ServiceInstance struct {
	ID        string   `json:"id"`        //注册中心的id
	Name      string   `json:"name"`      //注册中心name
	Version   string   `json:"version"`   //服务版本
	Endpoints []string `json:"endpoints"` //服务的地址 e.g：http://127.0.0.1:8000 grpc://127.0.0.1:9000
	Metadata  []string `json:"metadata"`  //服务元数据
}
