package znet

import (
	"errors"
	"fmt"
	"io"
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



	// 告知当前连接已经退出/停止 channel
	ExitChan chan bool

	// 无缓冲管道，用于读写 Goroutine 之间的通信
	msgChan chan []byte


	// 消息的管理 MsgID 和 对应的处理业务的API关系
	pMsgHandle ziface.IMsgHandle
}



// 初始化连接模块的方法
func NewConnection( pConn *net.TCPConn,iConnID uint32,msgHanle ziface.IMsgHandle) *Connection {
	pC := &Connection{
		pConn: pConn,
		iConnID: iConnID,
		pMsgHandle: msgHanle,
		isClosed: false,
		msgChan: make(chan []byte),
		ExitChan: make(chan bool,1),
	}
	return pC
}


// 当前连接的读数据的业务
func (pC *Connection)StartReadr()  {
	fmt.Println("[Reader Goroutine is runnig] ...")


	defer fmt.Println("iConnID=",pC.iConnID,"[Reader is exit]",pC.RemoteAddr().String())
	defer pC.Stop()

	for  {
		// 读取客户端的数据到buf 中
		//aBuf := make([]byte, utils.GlobalObject.IMaxPackageSize)
		//_,err  := pC.pConn.Read(aBuf)
		//if err != nil {
		//	fmt.Println("recv buf err=",err)
		//	continue
		//}
		// 创建一个拆包解包的对象
		objDP := NewDataPack()
		// 读取客户端的 Msg Head 二进制流 8 字节
		headData := make([]byte,objDP.GetHandLen())
		if _,err := io.ReadFull(pC.GetTCPConnection(),headData); err != nil {
			fmt.Println("read msg head error:",err)
			break
		}


		// 拆包，得到msgID 和 msgDatalen 放在 msg 消息中
		msg,err := objDP.Unpack(headData)
		if err != nil {
			fmt.Println("objDP Unpack  error:",err)
			break
		}
		// 根据msgDatalen 再次读取 data，放在 msg.data 消息中
		var aData[] byte
		if msg.GetMsgLen() > 0 {
			aData  = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(pC.GetTCPConnection(),aData);err != nil{
				fmt.Println("read msg data error:",err)
				break
			}
		}
		msg.SetMsgData(aData)

		// 得到 当前 conn 数据的 Request 请求
		objReq := Request{
			conn: pC,
			objMsg: msg,
		}

		// 从路由中，找到注册绑定的Conn 对应的router 调用
		// 根据绑定好的 msgID  找到对应处理业务的 API 业务执行
		go pC.pMsgHandle.DoMsgHandler(&objReq)
	}
}


// 提供一个sendmsg 方法， 将我们要发送给客户端的数据 ，先进性封包，再发送
func (pC *Connection)SendMsg(iMsgID uint32,aData []byte) error  {
	if pC.isClosed == true {
		return  errors.New("Connection is closed")
	}

	// 将data进行封包 msgdataLen | msgid | msgdata
	objDP := NewDataPack()
	binaryMsg,err := objDP.Pack(NewMsgPackage(iMsgID,aData))
	if err != nil{
		fmt.Println("Pack erro msg id=",iMsgID,"err=",err)
		return errors.New("Pack erro msg ")
	}

	// 将数据发送给 写协程 由写协程发送给	客户端
	pC.msgChan <- binaryMsg

	return nil
}



// 启动连接 让当前的连接准备开始工作
func  (pC *Connection)Start(){
	fmt.Println("conn start()..iConnID=",pC.iConnID)

	// 启动从当前连接的读数据的业务
	go pC.StartReadr()
	// 启动从当前连接写数据的业务
	go pC.StartWriter()

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

	// 告知 writer 关闭
	pC.ExitChan <- true

	// 关闭管道 回收资源
	close(pC.ExitChan)
	close(pC.msgChan)
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


// 写消息 goroutine , 专门发送给客户端消息的模块
func (pC *Connection)StartWriter()  {
	fmt.Println("[Writer Gortine is running ...]")
	defer fmt.Println(pC.RemoteAddr().String()," [conn Writer exit!]")
	// 不断阻塞的等待 channel 的消息， 进行发送消息到客户端
	for  {
		select {
		case data := <-pC.msgChan:
			// 有数据要写给客户端
			if _,err := pC.pConn.Write(data);err != nil{
				fmt.Println("send data err:",err)
				return
			}

		case <- pC.ExitChan:
			// 代表 Reader 已经退出，此时 writer 也要退出
			return
		}
	}

}






