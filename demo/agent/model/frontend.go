package model

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
)

var (
	FrontendListener cellnet.Peer
)

func GetClient(sesid int64) cellnet.Session {

	return FrontendListener.(peer.SessionManager).GetSession(sesid)
}
