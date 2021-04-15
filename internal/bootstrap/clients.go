package bootstrap

import "sync"

type ServiceRoutes struct {
	mutex        sync.RWMutex
	serviceAddrs map[string]string
}

func NewServiceRoutes() *ServiceRoutes {
	return &ServiceRoutes{
		serviceAddrs: make(map[string]string),
	}
}

func (r *ServiceRoutes) Get(serviceName string) (serviceAddr string, ok bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	serviceAddr, ok = r.serviceAddrs[serviceName]
	return
}

func (r *ServiceRoutes) Set(serviceName string, serviceAddr string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.serviceAddrs[serviceName] = serviceAddr
}
