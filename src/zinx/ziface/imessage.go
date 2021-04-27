package ziface

/*
	将请求的消息封装到一个 message 中，定义一个抽象的接口
*/
type IMessage interface {
	// 获取消息的ID
	GetMsgID() uint32
	// 获取消息的长度
	GetMsgLen() uint32
	// 获取消息的内容
	GetMsgData() []byte

	// 设置消息的ID
	SetMsgID(uint32)
	//设置消息的长度
	SetMsgLen(uint32)
	//设置消息的内容
	SetMsgData([]byte)
}







