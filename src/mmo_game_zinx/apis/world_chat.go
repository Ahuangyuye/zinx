package apis

import (
	"fmt"
	"pkg/mod/github.com/golang/protobuf/proto"
	"zinx/src/mmo_game_zinx/core"
	"zinx/src/mmo_game_zinx/pb/pb"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

// 世界聊天 路由业务

type WorldChatApi struct {
	znet.BaseRouter
}

func (pWc *WorldChatApi) Handle(request ziface.IRequest)  {
	// 1 解析客户端传递进来的 proto 协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(),proto_msg)
	if err != nil{
		fmt.Println("talk Unmarshal error ")
		return
	}

	// 2 当前的聊天信息是属于哪个玩家发送的
	pid,err := request.GetConnection().GetProperty("pid")

	// 3 根据pid得到对应的player对象 （pid 万能的指针）
	player :=  core.PWorldMgrObj.GetPlayerByPid(pid.(int32))

	// 4 将这个消息广播其他全部在线的玩家
	player.Talk(proto_msg.Content)
}













