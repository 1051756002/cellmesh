package service

import (
	"github.com/davyxu/cellnet"
)

type Event interface {
	// 给来源会话(网关,服务)发消息
	Session() cellnet.Session

	// 事件携带的消息
	Message() interface{}

	// 网关透传输出,如客户端在网关的SessionID
	PassThrough() interface{}

	// 回复客户端
	Reply(msg interface{})

	Raw() cellnet.Event
}

type EventFunc func(Event)

// 通用服务
type Service interface {
	// 服务发现注册
	Start()

	Stop()

	IsReady() bool
}

// 通讯用的服务
type CommunicateService interface {
	Service

	// 在cellnet中注册的事件处理器名
	SetProcessor(name string)

	// 接收消息的回调
	SetEventCallback(dis EventFunc)
}
