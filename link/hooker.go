package link

import (
	"github.com/davyxu/cellmesh"
	meshproto "github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
)

// 服务互联消息处理
type SvcEventHooker struct {
}

func (SvcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch msg := inputEvent.Message().(type) {
	case *meshproto.ServiceIdentifyACK: // 服务方收到连接方的服务标识

		if pre := GetLink(msg.SvcID); pre == nil {

			// 添加连接上来的对方服务
			markLink(inputEvent.Session(), msg.SvcID, msg.SvcName)
		}
	case *cellnet.SessionConnected:

		// 用Connector的名称（一般是ProcName）让远程知道自己是什么服务，用于网关等需要反向发送消息的标识
		inputEvent.Session().Send(&meshproto.ServiceIdentifyACK{
			SvcID:   cellmesh.GetLocalSvcID(),
			SvcName: cellmesh.ProcName,
		})

	case *cellnet.SessionClosed:

		if svcID := GetLinkSvcID(inputEvent.Session()); svcID != "" {
			log.SetColor("yellow").Infof("Remove service link : %s %s", GetLinkSvcID(inputEvent.Session()), getPeerDescString(inputEvent.Session().Peer()))
		}
	}

	return inputEvent

}

func (SvcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	// 服务器间通讯协议
	proc.RegisterProcessor("tcp.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(SvcEventHooker), new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	// 与客户端通信的处理器
	proc.RegisterProcessor("tcp.client", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
