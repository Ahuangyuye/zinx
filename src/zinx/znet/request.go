package znet

import "zinx/src/zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection

	// 客户端请求的数据
	//aData []byte
	objMsg ziface.IMessage
}


// 得到当前连接
func (pR *Request) GetConnection() ziface.IConnection {
	return  pR.conn
}

// 得到请求的消息数据
func(pR *Request)  GetData() []byte{
	return  pR.objMsg.GetMsgData()
}


// 得到请求的消息ID
func(pR *Request)  GetMsgID() uint32{
	return  pR.objMsg.GetMsgID()
}



