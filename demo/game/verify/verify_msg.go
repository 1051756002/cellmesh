package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
)

func Verify(req *proto.VerifyREQ, ack *proto.VerifyACK) {

	fmt.Printf("verfiy: %+v \n", req.GameToken)

	//proto.RouterBindUser()

	ack.Result = 0
}
