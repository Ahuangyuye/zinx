package znet

import (
	"fmt"
	"net"
	"zinx/src/zinx/ziface"
)

/*
连接模块
*/
type Connection struct {

	// 当前连接的socket TCP 套接字
	pConn *net.TCPConn

	// 连接的ID
	iConnID uint32

	// 当前的连接状态
	isClosed bool

	// 当前连接所绑定的处理业务方法API
	handleAPIfunc ziface.HandleFunc

	// 告知当前连接已经退出/停止 channel
	ExitChan chan bool


}



// 初始化连接模块的方法
func NewConnection( pConn *net.TCPConn,iConnID uint32,callback_api ziface.HandleFunc ) *Connection {
	pC := &Connection{
		pConn: pConn,
		iConnID: iConnID,
		handleAPIfunc: callback_api,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}
	return pC
}


// 当前连接的读数据的业务
func (pC *Connection)StartReadr()  {
	fmt.Println("Reader Coroutine is runnig ...")


	defer fmt.Println("iConnID=",pC.iConnID,"Reader is exit",pC.RemoteAddr().String())
	defer pC.Stop()

	for  {
		// 读取客户端的数据到buf 中
		aBuf := make([]byte,512)
		iCount,err  := pC.pConn.Read(aBuf)
		if err != nil {
			fmt.Println("recv buf err=",err)
			continue
		}
		// 调用当前连接所绑定的HandleAPI
		if err := pC.handleAPIfunc(pC.pConn,aBuf,iCount); err != nil {
			fmt.Println("iConnID=",pC.iConnID,"handleAPI is error ",err)
			break
		}
	}
}


// 启动连接 让当前的连接准备开始工作
func  (pC *Connection)Start(){
	fmt.Println("conn start()..iConnID=",pC.iConnID)

	// 启动从当前连接的读数据的业务
	go pC.StartReadr()
	// TODO 启动从当前连接写数据的业务


}

// 停止连接 结束当前连接的工作
func  (pC *Connection)Stop(){
	fmt.Println("conn stop .. iConnID=",pC.iConnID)

	// 如果当前连接已经关闭
	if pC.isClosed == true {
		return
	}

	pC.isClosed = true

	// 关闭socket 连接
	pC.pConn.Close()
	// 关闭管道 回收资源
	close(pC.ExitChan)
}


// 获取当前连接的绑定的 socket conn
func  (pC *Connection)GetTCPConnection() *net.TCPConn{
	return pC.pConn
}


// 获取当前连接模块的连接ID
func  (pC *Connection)GetConnID() uint32{
	return pC.iConnID
}


// 获取远程客户端的 TCP 状态的 ip  port
func  (pC *Connection)RemoteAddr() net.Addr{
	return pC.pConn.RemoteAddr()
}


// 发送数据 将数据发送给远程的客户端
func  (pC *Connection)Send(data []byte) error{
	return  nil
}


