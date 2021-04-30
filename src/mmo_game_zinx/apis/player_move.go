package apis

import (
	"fmt"
	"pkg/mod/github.com/golang/protobuf/proto"
	"zinx/src/mmo_game_zinx/core"
	"zinx/src/mmo_game_zinx/pb/pb"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

//  玩家移动路由 API
type MoveApi struct {
	znet.BaseRouter
}

func (pM *MoveApi)Handle(request ziface.IRequest)  {
	// 解析客户端传递进来的proto 协议
	proto_msg := &pb.Postion{}
	err := proto.Unmarshal(request.GetData(),proto_msg)
	if err != nil{
		fmt.Println("Move : Postion Unmarshal error ",err)
		return
	}

	// 得到当前发送位置的是哪个玩家
	pid,err :=request.GetConnection().GetProperty("pid")
	if err != nil{
		fmt.Println("GetProperty pid error ",err)
		return
	}
	fmt.Println("Player pid = %d ,move (%f,%f,%f,%f)\n",pid,proto_msg.X,proto_msg.Y,proto_msg.Z,proto_msg.V)

	// 给其他玩家进行当前玩家位置的信息广播
	player := core.PWorldMgrObj.GetPlayerByPid(pid.(int32))

	// 广播并更新当前玩家的坐标
	player.UnpdatePos(proto_msg.X,proto_msg.Y,proto_msg.Z,proto_msg.V)

}







