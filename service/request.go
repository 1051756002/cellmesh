package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"reflect"
)

func selectStrategy(descList []*discovery.ServiceDesc) *discovery.ServiceDesc {

	if len(descList) == 0 {
		return nil
	}

	return descList[0]
}

func Request(serviceName string, req interface{}, ackType reflect.Type, callback func(interface{})) error {

	addr, err := QueryServiceAddress(serviceName)
	if err != nil {
		return err
	}

	if rawConn, ok := connByAddr.Load(addr); ok {
		conn := rawConn.(Requestor)

		return conn.Request(req, ackType, callback)
	}

	return errors.New("connection not ready")
}
