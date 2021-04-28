package znet

import (
	"fmt"
	"net"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

// IServer 的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器的名称
	strName string

	// 服务器绑定的IP版本
	strIPVersion string

	// 服务器监听的iP
	strIP string

	// 服务器监听的端口
	iPort int

	// 当前的Server 添加一个 router , server 注册的连接对应的处理业务
	//objRouter ziface.IRouter

	// 当前 server 的消息管理模块，用了绑定 MsgID 和对应的处理业务API关系
	pMsgHandler ziface.IMsgHandle

	// 该server 的连接管理器
	pConnMgr ziface.IConnManager

	// 该 Server 创建连接之后自动调用 hook 函数
	OnConnStart func(conn ziface.IConnection)

	// 该 Server 断开连接之前自动调用 hook 函数
	OnConnStop func(conn ziface.IConnection)
}

//// 定义当前客户端连接的所绑定的handleAPI (目前这个handle 是写死的)，以后由调用者来	自定义
//func CallBackTOClient( pConn *net.TCPConn, aData[]byte,iDataLen int) error {
//	// 回写的业务
//	fmt.Println("[pConn Handle CallbackToClient]...")
//	if _,err := pConn.Write(aData[:iDataLen]); err != nil {
//		fmt.Println("write back buf err=",err)
//		return errors.New("CallBackTOClient error ")
//	}
//	return  nil
//}



// 启动服务器
func (s *Server) Start()  {

	fmt.Printf("[Zinx] server Name:%s listenner at IP:%s port:%d is starting \n", utils.GlobalObject.StrName,  utils.GlobalObject.StrHostIP,utils.GlobalObject.ITcpPort)

	fmt.Printf("[Zinx] Version=%s,MaxConn:%d,MaxPackageSize:%d\n",utils.GlobalObject.StrVersion,utils.GlobalObject.IMaxConn,utils.GlobalObject.IMaxPackageSize)

	//fmt.Printf("[start] server listenner at IP:%s port:%d is starting \n",s.strIP,s.iPort)

	go func() {

		// 0 开启消息队列及Worker工作池
		s.pMsgHandler.StartWorkerPool()


		// 1 获取一个TCP的Addr
		strIPandPort := fmt.Sprintf("%s:%d",s.strIP,s.iPort)
		strAddr,err := net.ResolveTCPAddr(s.strIPVersion,strIPandPort)
		if err != nil {
			fmt.Println("reslove tcp addr error:",err)
			return
		}

		// 2 监听服务器的地址
		pListenner,err  :=net.ListenTCP(s.strIPVersion,strAddr)
		if err != nil {
			fmt.Println("ListenTCP  strIPVersion = ",s.strIPVersion,"strAddr=",strAddr,"error:",err)
			return
		}
		fmt.Println("start Zinx server succ", s.strName,"succ listenning ...")

		var iCID  uint32
		iCID = 0
		// 3 阻塞的等待客户端链接，处理客户端的业务(读写)
		for {
			// 如果有客户端连接过来，阻塞会返回
			pConn,err := pListenner.AcceptTCP()
			if err != nil {
				fmt.Println(" AcceptTCP error:",err)
				continue
			}

			// 最大连接个数判断
			if s.pConnMgr.GetAllConnCount() >= utils.GlobalObject.IMaxConn {
				// TODO 给客户端相应的一个超出最大连接数的错误包
				fmt.Println("Connection beyond Max=",utils.GlobalObject.IMaxConn," Current connCount=",s.pConnMgr.GetAllConnCount() )
				pConn.Close()
				continue
			}

			// 将处理新连接的业务方法 和 conn 进行绑定，得到连接模块以便后续调用
			pDealConn  := NewConnection(s,pConn,iCID,s.pMsgHandler)
			iCID ++
			// 启动当前的连接业务处理
			go pDealConn.Start()
		}
	}()
}


// 停止服务器
func (s *Server) Stop()  {
	//  将一些服务器的资源，状态 或者一些已经开辟的连接信息，进行停止或者回收
	fmt.Println("[Stop Zinx Server] -->",s.strName)
	s.pConnMgr.ClearAllConn()
}

// 运行服务器
func (s *Server) Server()  {
	// 启动 server 的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外功能

	// 主线程阻塞状态
	select {

	}
}

// 路由功能，给当前的服务注册一个路由方法, 供客户端的连接处理使用
func (pS *Server)AddRouter(iMsgID uint32, router ziface.IRouter)  {
	pS.pMsgHandler.AddRouter(iMsgID,router)
	//pS.objRouter = router
	fmt.Println("add router succ !!!")
}

// 获取 connMgr
func (pS *Server)GetConnManager() ziface.IConnManager  {
	return pS.pConnMgr
}


// 初始化server 模块的方法
func NewServer(strNameIn string) ziface.IServer {
	//objS := &Server{
	//	strName:  strNameIn,
	//	strIPVersion: "tcp4",
	//	strIP: "0.0.0.0",
	//	iPort: 8999,
	//	objRouter: nil,
	//}
	objS := &Server{
		strName:  utils.GlobalObject.StrName, //utils.pGlobalObject,
		strIPVersion: "tcp4",
		strIP: utils.GlobalObject.StrHostIP,
		iPort: utils.GlobalObject.ITcpPort,
		//objRouter: nil,
		pMsgHandler:NewMsgHandler(),
		pConnMgr:NewConnManager(),
	}
	return objS
}


//// 该 Server 创建连接之后自动调用 hook 函数
//OnConnStart func(conn ziface.IConnection)
//
//// 该 Server 断开连接之前自动调用 hook 函数
//OnConnStop func(conn ziface.IConnection)

// 注册 OnConnStart 钩子函数的方法
func(pS *Server)SetOnConnStart(hookFunc func(connection ziface.IConnection )){
	pS.OnConnStart = hookFunc
}

// 注册 OnConnStop 钩子函数的方法
func(pS *Server)SetOnConnStop(hookFunc func(connection ziface.IConnection )){
	pS.OnConnStop = hookFunc
}


//调用 OnConnStart 钩子函数的方法
func(pS *Server)CallOnConnStart(connection ziface.IConnection ){
	if pS.OnConnStart != nil{
		fmt.Println("-- Call OnConnStart --")
		pS.OnConnStart(connection)
	}
}

// 调用 OnConnStop 钩子函数的方法
func(pS *Server)CallOnConnStop(connection ziface.IConnection ){
	if pS.OnConnStop != nil{
		fmt.Println("-- Call OnConnStop --")
		pS.OnConnStop(connection)
	}
}

