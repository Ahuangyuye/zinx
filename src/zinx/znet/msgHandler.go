package znet

import (
	"fmt"
	"strconv"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个 msgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter

	// 负责 worker 取任务的消息队列
	TaskQueueChan []chan ziface.IRequest

	// 业务工作池的worker 数量
	IWorkerPoolSize uint32
}

// 创建MsgHandle 方法
func NewMsgHandler() *MsgHandle {
	return &MsgHandle{
		Apis:            make(map[uint32]ziface.IRouter),
		IWorkerPoolSize: utils.GlobalObject.IWorkerPoolSize, // 从全局配置中获取
		TaskQueueChan:   make([]chan ziface.IRequest, utils.GlobalObject.IWorkerPoolSize),
	}
}

// 调度 执行 对于的router 消息处理方法
func (pMH *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从 request 获取 msgID
	handler, ok := pMH.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), " is not")
		return
	}

	// 2 调用对应的 router 业务
	handler.PreHandle(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

// 为消息添加具体的处理逻辑
func (pMH *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断当前 msg  绑定的 API 处理方法是否已经存在
	if _, ok := pMH.Apis[msgID]; ok {
		// id 已经注册了
		panic("repeat api ,msgID:" + strconv.Itoa(int(msgID)))
	}
	// 2 添加 msg 与 API 的 绑定关系
	pMH.Apis[msgID] = router
	fmt.Println("Add Api msgDI=", msgID, " succ")
}

// 启动一个worker 工作池(开启工作池的动作只能发生一次，一个zinx 框架只能由一个worker 工作池)
func (pMH *MsgHandle) StartWorkerPool() {
	// 根据  workerPoolSize 分别开启 worker,每个worker用一个 go 来承载
	for i := 0; i < int(pMH.IWorkerPoolSize); i++ {
		// 启动一个 worker
		// 1. 当前的worker 对应的channel 消息队列开辟空间
		pMH.TaskQueueChan[i] = make(chan ziface.IRequest,utils.GlobalObject.IMaxWorkerTaskLen)
		// 2 启动当前的 worker ，阻塞等待消息从 channel 传递过来
		go pMH.StartOneWorker(i,pMH.TaskQueueChan[i] )
	}

}

// 启动一个worker 工作流程
func (pMH *MsgHandle) StartOneWorker(iWorkerID int,taskQueueChan chan ziface.IRequest) {
	fmt.Println("worker ID:",iWorkerID,"is started ...")

	// 不断的阻塞等待对应的消息队列的消息
	for {
		select {
			// 如果有消息过来，出列的就是一个客户端的 Request， 执行的是当前的 Request 所绑定的业务
			case request := <- taskQueueChan:
				pMH.DoMsgHandler(request)
		}
	}
}


// 将消息发送给消息任务队列 TaskQueue, 由 worker 进行处理
func (pMH *MsgHandle)SendMsgToTaskQueue(request ziface.IRequest)  {
	// 1 将消息平均分配给不通过的 worker
	// 根据客户端建立的ConnID进行分配
	iWorkerID := request.GetConnection().GetConnID() % pMH.IWorkerPoolSize
	fmt.Println("Add ConnID:",request.GetConnection().GetConnID()," request msgID:",request.GetMsgID()," to workerID:",iWorkerID)

	// 2 将消息发送给对应的worker的TaskQueue 即可
	pMH.TaskQueueChan[iWorkerID] <- request
}














