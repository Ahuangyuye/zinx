package ziface




/*
	封包、拆包 模块
	直接面向TCP 连接的数据流，用于处理 TCP 粘包问题
*/

type IDataPack interface {
	//  获取包的头的长度的方法
	GetHandLen() uint32

	// 封包方法
	Pack(msg IMessage)([]byte,error)
	// 拆包方法
	Unpack([]byte)(IMessage,error)
}
