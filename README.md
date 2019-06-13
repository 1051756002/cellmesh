 [![Build Status][3]][4] [![Go Report Card][5]][6] [![MIT licensed][11]][12] [![GoDoc][1]][2]

[1]: https://godoc.org/github.com/davyxu/cellmesh?status.svg
[2]: https://godoc.org/github.com/davyxu/cellmesh
[3]: https://travis-ci.org/davyxu/cellmesh.svg?branch=master
[4]: https://travis-ci.org/davyxu/cellmesh
[5]: https://goreportcard.com/badge/github.com/davyxu/cellmesh
[6]: https://goreportcard.com/report/github.com/davyxu/cellmesh
[11]: https://img.shields.io/badge/license-MIT-blue.svg
[12]: LICENSE

# cellmesh
基于cellnet的游戏服务框架

# 特点

## Based On Service Discovery(基于服务发现)

   通过服务发现自动实现服务互联,探测,挂接.无需配置服务器间的端口.

## Zero Config File(零配置文件)

   任何服务器配置均通过服务发现的KV存取,无任何配置文件.

## Code Generation(代码生成)

   基于github.com/davyxu/protoplus的代码生成技术,迅速粘合逻辑与底层,代码好看易懂且高效.

   使用更简单更强大的schema编写协议, 并自动生成protobuf schema.

## Transport based on cellnet(基于cellnet的网络传输)

   提供强大的扩展及适配能力.

# 使用cellmesh

    cellmesh 使用go module管理源码依赖， 所以确保go版本在1.12以上

## 下载cellmesh源码

```
    go get github.com/davyxu/cellmesh
```

# 概念

## Service（服务）

一个Service为一套连接器或侦听器，挂接消息处理的派发器Dispatcher

- 侦听端口自动分配

  Service默认启动时以地址:0启动，网络底层自动分配端口，由cellmesh将服务信息报告到服务发现

  其他Service发现新的服务进入网络时，根据需要自动连接服务

## Connection Management（连接维护）

从服务发现的服务信息，创建到不同服务间的长连接。同时在这些连接断开时维护连接

逻辑中根据策略从已有连接及信息选择出连接与目标通信，例如：选择负载最低的一台游戏服务器

# 目录结构
```
discovery
   kvconfig
      配置的快速获取接口。
   memsd
      面向游戏优化的服务发现实现。
service
   服务通信基础，以及服务发现封装。
shell
   框架通用的shell脚本。
tools
   protogen
      协议生成器，生成Go的消息绑定以及消息响应入口。
   routegen
      路由配置生成器。生成的配置可由agent动态读取并更新路由规则。
util
   所有框架通用的工具代码。

```

# 服务参数
 service包为服务进程提供命令行参数支持。服务进程的命令行参数同时也可以使用FlagFile像使用配置文件一样批量设置进程配置(参考demo/cfg/LocalFlag.cfg)

 详细参数说明如下：

- sdaddr

   指定服务发现服务器(memsd)地址, 通过服务发现,服务器可以快速获取配置以及其他可连接服务器地址,实现服务互联.

- svcgroup

   指定服务器分组. 一般情况下,认为一台物理机归属于一个svcgroup. 当然,也可以在一台物理机上放置多个分组,比如开发阶段.

   服务器分组也能方便服务器打包以及定位服务器位置.

- svcindex

   指定服务器索引, 标识同类服务器的多个不同进程,同类中的svcindex必须唯一,逻辑上,svcindex还会与uuid关联.

- wanip

   指定服务器所在物理机的外网IP,方便通知客户端要连接的IP,例如:login通知客户端game的外网IP.

- logcolor

   对日志着色, 规则参见github/davyxu/golog中的color.go, 写入文件时,请关闭此选项,避免日志中出现着色字符.

- logfile

   将日志写入文件,并不再输出到控制台.

- logfilesize

   指定输出日志文件单个大小,可使用B/M/G标识, 注意: golog默认不是rolling方式,日志会写入到尾数最大的日志文件.

- loglevel

   指定日志输出级别, 格式 日志名|级别, 日志名支持正则表达式, 级别可以为error, info等

- mutemsglog

   屏蔽指定消息的日志,多个消息使用'|'分割

- flagfile

   使用FlagFile格式(参考demo/cfg/LocalFlag.cfg),作为进程的命令行参数

# demo工程架构
![arch](demo/doc/architecture.png)

- login

    登录服, 短连接获取服务器列表及平台验证, 通过JWT生成token交由客户端

- game

    养成服, 网关后负责养成逻辑,使用客户端上传的JWT token验证客户端身份

- agent

    网关, 负责与客户端保持长连接,通过心跳处理客户端连接生命期, 客户端消息转发

- hub

    中转, 处理服务器人数负载, 使用模式为订阅分发模式,处理跨服通知

- client

    模拟客户端, 模拟客户端逻辑


# 开发进度
- [x] 基于Consul的服务发现及Watch机制
- [x] 网关基本逻辑
- [x] 带服务发现的连接器,侦听器
- [x] 网关会话绑定
- [x] 服务发现Consul的KV封装
- [x] 网关路由规则生成,上传和下载
- [x] 系统消息响应入口
- [x] 网关广播
- [x] 网关心跳处理
- [x] 频道订阅发布hub
- [x] 服务在线人数更新及连接选择
- [x] 支持WebSocket网关及登录协议
- [ ] 登录服务器，JWT验证
- [ ] 玩家数据读取
- [ ] 游戏服务器，成长逻辑，花钱升级等级
- [ ] 社交服务器，聊天逻辑
- [ ] 机器人

# 参考及使用建议
由于本框架尚在开发中,demo在不断完善添加新功能. 因此响应的接口和设计会经常性发生变化.


# Tips
## 为什么使用memsd的服务发现替换consul？
早期版本的cellmesh使用consul作为服务发现，cellmesh使用主动汇报服务信息的方式保证consul中能及时更新服务信息。
但实际使用中发现有如下问题：
1. 偶尔出现高CPU占用，Windows休眠恢复后也会造成严重的高CPU现象。
2. consul的API并没有本地cache，需要高速查询时，并没有很好的性能。
3. 多服更新时没有原子更新，容易形成严重的不同步现象。
4. 依赖重，代码量巨大，使用vendor而不是go module方式管理代码，编译慢。
基于以上考虑，决定兼容服务发现接口，同时编写对游戏服务友好的发现系统：memsd。



友情提示：
    demo工程仅是随框架附带的实例代码，建议将demo为蓝本建立自己的工程蓝本。


# 备注

感觉不错请star, 谢谢!

开源讨论群: 527430600 验证请发cellmesh

知乎: http://www.zhihu.com/people/sunicdavy

提交bug及特性: https://github.com/davyxu/cellmesh/issues
