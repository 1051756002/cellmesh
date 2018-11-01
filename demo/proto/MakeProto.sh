#!/usr/bin/env bash

CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`

set -e
Protoc=${GOPATH}/bin/protoc

CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/tools/protogen

ProtoPlusGen=${GOPATH}/bin/protoplus
go build -v -o ${ProtoPlusGen} github.com/davyxu/protoplus

RouteGen=${GOPATH}/bin/routegen
go build -v -o ${RouteGen} github.com/davyxu/cellmesh/tools/routegen

cd ${CURRDIR}

# windows下时，添加后缀名
if [ `go env GOHOSTOS` == "windows" ];then
	EXESUFFIX=.exe
fi

echo "生成服务器协议的go消息..."
${ProtoPlusGen} -package=proto -go_out=msgsvc_gen.go `source ./proto_svc.sh`

echo "生成服务器协议的消息绑定..."
${CellMeshProtoGen} -package=proto -cmgo_out=msgbind_gen.go `source ./proto_all.sh`


echo "生成客户端协议的protobuf proto文件..."
${ProtoPlusGen} --package=proto -pb_out=clientmsg_gen.proto `source ./proto_client.sh`

echo "生成客户端协议的protobuf的go消息...."
${Protoc} --plugin=protoc-gen-gogofaster=${GOPATH}/bin/protoc-gen-gogofaster${EXESUFFIX} --gogofaster_out=. --proto_path=. clientmsg_gen.proto
rm -f ./clientmsg_gen.proto


echo "更新agent路由表"
${RouteGen} -configpath=config_demo/route_rule `source ./proto_client.sh`
