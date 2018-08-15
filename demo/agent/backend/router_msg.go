package backend

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func RouterBindUser(ev *service.Event, req *proto.RouterBindUserREQ, ack *proto.RouterBindUserACK) {

	fmt.Println("bind user", req.Token)
}
