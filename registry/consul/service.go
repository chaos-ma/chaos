package consul

/**
* created by mengqi on 2023/11/13
 */

import (
	"sync"
	"sync/atomic"

	"chaos/registry"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    *atomic.Value
	lock        sync.RWMutex
}

func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	//原子操作， 保证线程安全， 我们平时写struct的时候
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}
