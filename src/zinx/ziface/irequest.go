package ziface

/*
IRequest接口：
实际上是把客户端请求的连接信息 和 请求的数据 包装到一个 request 中
*/

type IRequest interface {

	// 得到当前连接
	GetConnection() IConnection

	// 得到请求的消息数据
	GetData() []byte

	// 得到请求的消息ID
	GetMsgID() uint32
}


