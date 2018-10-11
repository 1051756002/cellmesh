package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"time"
)

func (self *consulDiscovery) startRefresh() {

	for {

		time.Sleep(time.Second * 3)

		svc, err := self.client.Agent().Services()
		if err != nil {
			log.Errorln(err)
			continue
		}

		var newList []*discovery.ServiceDesc

		for _, detail := range svc {
			newList = append(newList, consulSvcToService(detail))
		}

		self.newCacheGuard.Lock()
		self.newCache = newList
		self.newCacheGuard.Unlock()

		self.OnCacheUpdated("add", nil)
	}

}

func (self *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {

	//
	//if raw, ok := self.cache.Load(name); ok {
	//	ret = raw.([]*discovery.ServiceDesc)
	//}
	self.newCacheGuard.RLock()
	for _, sd := range self.newCache {
		if sd.Name == name {
			ret = append(ret, sd)
		}
	}
	self.newCacheGuard.RUnlock()

	return
}

func (self *consulDiscovery) QueryAll() (ret []*discovery.ServiceDesc) {

	self.newCacheGuard.RLock()
	copy(ret, self.newCache)
	self.newCacheGuard.RUnlock()

	//self.cache.Range(func(key, value interface{}) bool {
	//	ret = append(ret, value.([]*discovery.ServiceDesc)...)
	//
	//	return true
	//})

	return
}

// from github.com/micro/go-micro/registry/consul_registry.go
func (self *consulDiscovery) directQuery(name string) (ret []*discovery.ServiceDesc, err error) {

	result, _, err := self.client.Health().Service(name, "", false, nil)

	if err != nil {
		return nil, err
	}

	for _, s := range result {

		if s.Service.Service != name {
			continue
		}

		if isServiceHealth(s) {

			sd := consulSvcToService(s.Service)

			log.Debugf("  got servcie, %s", sd.String())

			ret = append(ret, sd)
		}

	}

	return

}

func (self *consulDiscovery) RegisterNotify(mode string) (ret chan struct{}) {

	ret = make(chan struct{})

	self.notifyGuard.Lock()
	switch mode {
	case "add":
		self.addNotify = append(self.addNotify, ret)
	case "remove":
		self.removeNotify = append(self.removeNotify, ret)
	}
	self.notifyGuard.Unlock()

	return
}

func (self *consulDiscovery) DeregisterNotify(mode string, c chan struct{}) {

	self.notifyGuard.Lock()
	switch mode {
	case "add":
		for index, n := range self.addNotify {
			if n == c {
				self.addNotify = append(self.addNotify[:index], self.addNotify[index+1:]...)
				break
			}
		}
	case "remove":
		for index, n := range self.removeNotify {
			if n == c {
				self.removeNotify = append(self.removeNotify[:index], self.removeNotify[index+1:]...)
				break
			}
		}
	}
	self.notifyGuard.Unlock()

}

func (self *consulDiscovery) OnCacheUpdated(eventName string, desc *discovery.ServiceDesc) {

	self.notifyGuard.RLock()
	switch eventName {
	case "add":
		//log.Debugf("Add service '%s'", desc.ID)

		for _, n := range self.addNotify {
			n <- struct{}{}
		}

	case "remove":
		//log.Debugf("Remove service '%s'", desc.ID)

		for _, n := range self.removeNotify {
			n <- struct{}{}
		}
	}

	self.notifyGuard.RUnlock()
}
