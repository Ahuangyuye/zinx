package znet

type Message struct {
	IID      uint32 // 消息的ID
	IDataLen uint32 // 消息的长度
	AData    []byte // 消息的内容
}

// 创建一个 message  消息包
func NewMsgPackage(iID uint32,aData []byte) *Message  {

	return &Message{
		IID: iID,
		IDataLen: uint32(len(aData)),
		AData: aData,
	}
}



// 获取消息的ID
func (pM *Message)GetMsgID() uint32{
	return pM.IID
}
// 获取消息的长度
func (pM *Message)GetMsgLen() uint32{
	return pM.IDataLen
}
// 获取消息的内容
func (pM *Message)GetMsgData() []byte{
	return pM.AData
}

// 设置消息的ID
func (pM *Message)SetMsgID(iid uint32){
	pM.IID = iid
}
//设置消息的长度
func (pM *Message)SetMsgLen(iLen uint32){
	pM.IDataLen = iLen
}
//设置消息的内容
func (pM *Message)SetMsgData(aData []byte){
	pM.AData = aData
}


