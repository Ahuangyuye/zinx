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

	err := request.GetConnection().SendMsg(200,[]byte("hello world... "))
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






// ping test 自定义路由
type HelloZinxRouter struct {
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
func (this *HelloZinxRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handler...")
	//_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping... ping...ping...ping...\n"))
	//if err != nil {
	//	fmt.Println("call back Handler err=",err)
	//}
	fmt.Println("recv from client--> msgID:",request.GetMsgID(),
		" msgData:",string(request.GetData()))

	err := request.GetConnection().SendMsg(201,[]byte("HelloZinxRouter world... "))
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


// 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection)  {
	fmt.Println("=========== DoConnectionBegin ===========")
	if err:=conn.SendMsg(100,[]byte("DOCONNECTIONBEGIN")) ; err != nil {
		fmt.Println(err)
	}

	// 给当前连接设置一些属性
	fmt.Println("set conn Name , hoe ...")
	conn.SetProperty("name","xiongmao")
	conn.SetProperty("age","18")
	id := conn.GetConnID()
	conn.SetProperty("id",id)
}


// stop连接之前执行的钩子函数
func DoConnectionLost(conn ziface.IConnection)  {
	fmt.Println("=========== DoConnectionLost ===========")

	// 获取连接属性
	if name,err := conn.GetProperty("name");err == nil {
		fmt.Println("name-->",name)
	}
	if id,err := conn.GetProperty("id");err == nil {
		fmt.Println("id -->",id)
	}
}


func main()  {
	// 1 创建一个server 句柄，使用 zinx 的api
	pS := znet.NewServer("[zinx V0.7]")

	// 2 注册连接 hook 钩子函数
	pS.SetOnConnStart(DoConnectionBegin)
	pS.SetOnConnStop( DoConnectionLost)
	// 3 给当前 zinx 框架添加一个 自定义的 router
	pS.AddRouter(0,&PingRouter{})
	pS.AddRouter(1,&HelloZinxRouter{})

	// 4 启动server
	pS.Server()

}


