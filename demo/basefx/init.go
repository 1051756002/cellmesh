package basefx

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
)

// 初始化框架
func Init(procName string) {

	msglog.BlockMessageLog("proto.PingACK")
	msglog.BlockMessageLog("proto.SvcStatusACK")

	fxmodel.Queue = cellnet.NewEventQueue()

	fxmodel.Queue.StartLoop()

	service.Init(procName)
}

// 等待退出信号
func StartLoop(onReady func()) {

	fxmodel.CheckReady()

	if onReady != nil {
		cellnet.QueuedCall(fxmodel.Queue, onReady)
	}

	service.WaitExitSignal()
}

// 退出处理
func Exit() {
	fxmodel.StopAllService()
}
