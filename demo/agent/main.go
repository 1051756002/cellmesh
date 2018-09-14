package main

import (
	_ "github.com/davyxu/cellmesh/demo/agent/backend"
	"github.com/davyxu/cellmesh/demo/agent/frontend"
	"github.com/davyxu/cellmesh/demo/agent/heartbeat"
	"github.com/davyxu/cellmesh/demo/agent/routerule"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	"github.com/davyxu/cellmesh/service/cellsvc"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	routerule.Download()

	heartbeat.StartCheck()

	acc := cellsvc.NewAcceptor("svc->agent")
	acc.SetProcessor("tcp.ltv")
	acc.SetEventCallback(proto.GetDispatcher("router"))
	acc.Start()

	frontend.Start()

	util.WaitExit()

	frontend.Stop()
	acc.Stop()
}
