package main

import (
	"fmt"
	"zinx/src/mmo_game_zinx/apis"
	"zinx/src/mmo_game_zinx/core"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

// 当前客户端建立连接之后的 hook 函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个 player 对象
	player := core.NewPlayer(conn)

	// 给客户端发送 MsgID：1 消息
	player.SyncPid()

	// 给客户端发送msgID : 200 消息
	player.BroadCastStartPostion()


	// 将上线的玩家添加到 WorldManager 中
	core.PWorldMgrObj.AddPlayer(player)

	// 当前新上线的玩家 连接绑定一个PID
	conn.SetProperty("pid",player.Pid)

	// 同步周边玩家，告知他们当前玩家已经上线，广播当前玩家的位置信息
	player.SyncSurrounding()


	fmt.Println("===> Player Pid=", player.Pid)
}







func main() {

	pServer := znet.NewServer("MMO Game Zinx")

	// 连接和销毁 hook 钩子函数
	pServer.SetOnConnStart(OnConnectionAdd)

	// 注册一些路由业务
	pServer.AddRouter(2,&apis.WorldChatApi{})

	// 启动服务
	pServer.Server()

}
