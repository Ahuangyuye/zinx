package znet

import (
	"fmt"
	"net"
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
}
// 启动服务器
func (s *Server) Start()  {
	go func() {
		fmt.Printf("[start] server listenner at IP:%s port:%d is starting \n",s.strIP,s.iPort)
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

		// 3 阻塞的等待客户端链接，处理客户端的业务(读写)
		for {
			// 如果有客户端连接过来，阻塞会返回
			pConn,err := pListenner.AcceptTCP()
			if err != nil {
				fmt.Println(" AcceptTCP error:",err)
				continue
			}
			// 已经与客户端建立连接，做业务
			go func() {
				for {
					bBuf := make([]byte,512)
					iCount,err := pConn.Read(bBuf)
					if err != nil {
						fmt.Println("  pConn.Read(bBuf) error:",err)
						continue
					}

					// 回写
					if _,err := pConn.Write(bBuf[:iCount]); err != nil {
						fmt.Println("  pConn.Write(bBuf) error:",err)
						continue
					}
				}
			}()

		}
	}()
}


// 停止服务器
func (s *Server) Stop()  {
	// TODO 将一些服务器的资源，状态 或者一些已经开辟的连接信息，进行停止或者回收


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

// 初始化server 模块的方法
func NewServer(strNameIn string) ziface.IServer {
	objS := &Server{
		strName:  strNameIn,
		strIPVersion: "tcp4",
		strIP: "0.0.0.0",
		iPort: 8999,
	}
	return objS
}





