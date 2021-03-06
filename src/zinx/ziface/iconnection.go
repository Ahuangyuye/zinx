package ziface

import "net"

// 定义连接模块的抽象层
type IConnection interface {
	// 启动连接 让当前的连接准备开始工作
	Start()

	// 停止连接 结束当前连接的工作
	Stop()

	// 获取当前连接的绑定的 socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的连接ID
	GetConnID() uint32

	// 获取远程客户端的 TCP 状态的 ip  port
	RemoteAddr() net.Addr

	// 发送数据 将数据发送给远程的客户端
	SendMsg(iMsgID uint32,aData []byte) error

	//设置连接属性
	SetProperty(strKey string,value interface{})
	//获取链接属性
	GetProperty(strKey string) (interface{},error)
	//移除连接属性
	RemoveProperty(strKey string)
}


// 定义一个处理连接业务的方法
// 三个形参，一个返回值为error 的类型
// [*net.TCPConn:连接句柄] [ []byte: 具体数据内容 ] [int:数据长度]
type  HandleFunc func( *net.TCPConn, []byte, int) error


