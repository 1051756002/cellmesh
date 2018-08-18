package heartbeat

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet/timer"
	"time"
)

func StartCheck() {

	// 从KV获取配置,默认关闭
	heatBeatDuration := kvconfig.Int32("config/agent/heatbeat_sec", 0)

	// 为0时关闭心跳检查
	if heatBeatDuration == 0 {
		return
	}

	// 接收客户端心跳
	proto.Handle_Agent_PingACK = func(ev service.Event) {
		u := model.SessionToUser(ev.Session())
		if u != nil {
			u.LastPingTime = time.Now()
		}
	}

	// 超时检查比心跳稍长
	timeOutDur := time.Duration(heatBeatDuration+5) * time.Second

	log.Infof("Heatbeat duration: '%ds' ", heatBeatDuration)

	// 心跳检查
	timer.NewLoop(nil, timeOutDur, func(loop *timer.Loop) {

		now := time.Now()

		model.VisitUser(func(u *model.User) bool {

			if now.Sub(u.LastPingTime) > timeOutDur {
				log.Warnf("Close client due to heatbeat time out, id: %d", u.ClientSession.ID())
				u.ClientSession.Close()
			}

			return true
		})

	}, nil).Start()
}
