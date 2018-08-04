package consulsd

import (
	"context"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"time"
)

// 本地服务更新TTL
type localService struct {
	Desc   *discovery.ServiceDesc
	Cancel context.CancelFunc

	ctx context.Context

	agent *api.Agent
}

const healthWords = "cellmesh service ready"

func (self *localService) Update() {

	//log.Debugf("UpdateTTL id: %s begin", self.ID)

	for {

		select {
		case <-self.ctx.Done():
			return
		default:

			//log.Debugf("UpdateTTL id: %s", self.ID)

			self.agent.UpdateTTL(self.Desc.ID, healthWords, "pass")

			time.Sleep(ServiceTTL)
		}
	}

	//log.Debugf("UpdateTTL id: %s end", self.ID)
}

func (self *localService) Stop() {
	self.Cancel()
}

func newLocalService(svc *discovery.ServiceDesc, agent *api.Agent) *localService {

	ctx, cancel := context.WithCancel(context.Background())

	self := &localService{
		Desc:   svc,
		Cancel: cancel,
		ctx:    ctx,
		agent:  agent,
	}

	go self.Update()

	return self
}
