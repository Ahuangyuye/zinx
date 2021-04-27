package main

import (
	"fmt"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

/*
基于Zinx 框架来开发的 服务器端应用程序

*/


// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}
//// Test PreRouter
//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	// request.GetConnection() :获取包装好的 conn
//	// request.GetConnection().GetTCPConnection() ： 获取原生的tcpconn
//	// Write 写数据
//	_,err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("call back PreHandle err=",err)
//	}
//}

// Test PreRouter
func (this *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call Router Handler...")
	//_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping... ping...ping...ping...\n"))
	//if err != nil {
	//	fmt.Println("call back Handler err=",err)
	//}
	fmt.Println("recv from client--> msgID:",request.GetMsgID(),
		" msgData:",string(request.GetData()))

	err := request.GetConnection().SendMsg(1,[]byte("hello world... "))
	if err != nil{
		fmt.Println(err)
	}

}
//// Test PreRouter
//func (this *PingRouter) PostHandler(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandler...")
//	_,err := request.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
//	if err != nil {
//		fmt.Println("call back PostHandler err=",err)
//	}
//}



func main()  {
	// 1 创建一个server 句柄，使用 zinx 的api
	pS := znet.NewServer("[zinx V0.5]")

	// 2 给当前 zinx 框架添加一个 自定义的 router
	pS.AddRouter(&PingRouter{})

	// 3 启动server
	pS.Server()

}


