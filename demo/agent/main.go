package main

import (
	"github.com/davyxu/cellmesh/demo/agent/router"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
)

func main() {

	svcfx.Init()

	go service.PrepareConnection("demo.game", service.NewRPCRequestor, nil)

	router.Start()

	util.WaitExit()

	router.Stop()
}
