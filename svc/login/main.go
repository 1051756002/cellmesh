package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	cellmesh.Init("login")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	// 服务互联
	link.LinkService(&link.ServiceParameter{
		SvcName: "hub",
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
