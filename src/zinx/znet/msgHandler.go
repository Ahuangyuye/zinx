package znet

import (
	"fmt"
	"strconv"
	"zinx/src/zinx/ziface"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个 msgID 所对应的处理方法
	Apis map[uint32] ziface.IRouter
}

// 创建MsgHandle 方法
func NewMsgHandler() *MsgHandle {
	return &MsgHandle{
		Apis: make( map[uint32]ziface.IRouter),
	}
}

// 调度 执行 对于的router 消息处理方法
func(pMH *MsgHandle)DoMsgHandler(request ziface.IRequest){
	// 1 从 request 获取 msgID
	handler,ok := pMH.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=",request.GetMsgID()," is not")
	}

	// 2 调用对应的 router 业务
	handler.PreHandle(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

// 为消息添加具体的处理逻辑
func(pMH *MsgHandle)AddRouter(msgID uint32, router ziface.IRouter){
	// 1 判断当前 msg  绑定的 API 处理方法是否已经存在
	if _,ok := pMH.Apis[msgID];ok{
		// id 已经注册了
		panic("repeat api ,msgID:" + strconv.Itoa(int(msgID)))
	}
	// 2 添加 msg 与 API 的 绑定关系
	pMH.Apis[msgID] = router
	fmt.Println("Add Api msgDI=",msgID," succ")
}







