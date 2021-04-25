package main

import (
	"zinx/src/zinx/znet"
)

/*
基于Zinx 框架来开发的 服务器端应用程序

*/

func main()  {
	// 1 创建一个server 句柄，使用 zinx 的api
	pS := znet.NewServer("[zinx V0.1]")

	// 2 启动server
	pS.Server()

}


