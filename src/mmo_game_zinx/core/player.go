package core

import (
	"fmt"
	"math/rand"
	"pkg/mod/github.com/golang/protobuf/proto"
	"sync"
	"zinx/src/mmo_game_zinx/pb/pb"
	"zinx/src/zinx/ziface"
)

/*
	玩家
*/

type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 当前玩家的连接（与客户端进行通信的 conn)

	X float32 // 平面 X 坐标
	Y float32 // 高度
	Z float32 // 平面 Y 坐标（不是Y）
	V float32 // 旋转的角度 0-360
}

// Player ID 生成器
var PidGen int32 = 1
var IDLock sync.Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	IDLock.Lock()
	id := PidGen
	PidGen++
	IDLock.Unlock()

	// 创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
	return p
}

// 提供一个发送给客户端消息的方法
// 主要是将 PB 的 protobuf 数据序列化之后，再调用 zinx 的 sendMsg 方法
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 将 proto message 结构体序列化，转换二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}

	// 将二进制文件 通过 zinx 框架的sendMsg 将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println(" player sendMsg error !!!", err)
		return
	}
	return
}

// 告知客户端玩家 pid ，同步已经生成的 玩家ID给客户端
func (p *Player) SyncPid() {
	// 组建MsgID：0 的 proto 数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPostion() {
	// 组建MsgID：200 的 proto 数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, // Tp2 代表广播位置的坐标
		Data: &pb.BroadCast_P{
			P: &pb.Postion{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// 玩家广播世界聊天消息
func (p *Player) Talk(strConent string) {
	// 1 组件 MsgID:200 proto 数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, // 世界聊天 type
		Data: &pb.BroadCast_Content{
			Content: strConent,
		},
	}

	// 2 得到当前世界所有的在线玩家
	players := PWorldMgrObj.GetAllPlayers()

	// 3 向所有的玩家（包括自己） 发送 msgID：200 的消息
	for _, player := range players {
		// player 分别给对应的 客户端发送 msg
		player.SendMsg(200, proto_msg)
	}
}

// 同步周边玩家，告知他们当前玩家已经上线
func (p *Player) SyncSurrounding(strConent string) {
	// 1 获取当前玩家周围有哪些玩家（九宫格）
	pids := PWorldMgrObj.PAoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, PWorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	// 2 将当前玩家的位置信息通过 MsgID：200 发给周围的玩家（让其他玩家看到自己）
	// 2.1 组建 MsgID：200 proto 数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, // 广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Postion{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.2 分别给周围的玩家发送 200 的消息 proto_msg
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	// 3 将周围的全部玩家的位置信息发给当前的玩家客户端 msgID:202（让自己看到其他的玩家）
	// 3.1 组建 MsgID：202 的proto 数据
	// 3.1.1 制作 pb.Player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		// 制作一个  msg player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Postion{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg,p)
	}
	// 3.1.2 封装 SyncPlayers proto 数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		//Ps: players_proto_msg,
		Ps: players_proto_msg[:],
	}

	// 3.2 将组建好的数据发送给当前客户端
	p.SendMsg(202,SyncPlayers_proto_msg)
}


//广播并更新当前玩家的坐标
func (p *Player) UnpdatePos(x,y,z,v float32){
	// 更新当前玩家的 player 对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	// 组建广播 proto 协议 MsgID:200 Tp=4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 4,	// 4 - 移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P:&pb.Postion{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取当前玩家的周边玩家 AOI 九宫格之内的玩家
	player := p.GetSuroundingPlayers()

	// 一次给每个玩家对应的客户端发送当前玩家位置更新的消息


}



//	// 获取当前玩家的周边玩家 AOI 九宫格之内的玩家
func (p *Player) GetSuroundingPlayers()[]*Player {

	pids := PWorldMgrObj.PAoiMgr.GetPidsByPos(p.X,p.Z)

	players := make([]*Player,0,len(pids))
	for _,pid := range  pids{
		players = append(players)
	}
}
