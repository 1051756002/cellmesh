package service

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
	"time"
)

type Requestor interface {
	Request(req interface{}, ackType reflect.Type, callback func(interface{})) error

	Session() cellnet.Session

	Start()

	IsReady() bool

	Stop()
}

var (
	connByAddr sync.Map
)

func GetSession(addr string) cellnet.Session {
	if rawConn, ok := connByAddr.Load(addr); ok {
		conn := rawConn.(Requestor)

		return conn.Session()
	}

	return nil
}

func GetRequestor(addr string) Requestor {
	if rawConn, ok := connByAddr.Load(addr); ok {
		return rawConn.(Requestor)
	}

	return nil
}

func QueryServiceAddress(serviceName string) (string, error) {
	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return "", err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return "", errors.New("target not reachable")
	}

	return fmt.Sprintf("%s:%d", desc.Address, desc.Port), nil
}

func connLoop(serviceName string) {

	for {
		addr, err := QueryServiceAddress(serviceName)

		if err == nil {
			readyConn := make(chan string)

			requestor := NewRPCRequestor(addr, readyConn)

			requestor.Start()

			if requestor.IsReady() {

				connByAddr.Store(addr, requestor)
				log.SetColor("green").Debugln("service ready: ", serviceName)

				// 连接断开
				<-readyConn
				connByAddr.Delete(addr)

				log.SetColor("yellow").Debugln("service invalid: ", serviceName)
			} else {

				requestor.Stop()
				time.Sleep(time.Second * 3)
				continue
			}

			requestor.Stop()
		}

		discovery.Default.WaitAdded()
	}
}

func PrepareConnection(serviceName string) {

	go connLoop(serviceName)
}

func WaitConnectionReady(serviceName string) {

	for {
		addr, err := QueryServiceAddress(serviceName)

		if err != nil {
			continue
		}

		if rawConn, ok := connByAddr.Load(addr); ok {
			conn := rawConn.(Requestor)

			if conn.IsReady() {
				return
			}
		}

	}

}
