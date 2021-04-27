package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

// 封包 拆包的具体模块

type DataPack struct {

}

// 封包 拆包 实例的一个初始化方法
func NewDataPack() *DataPack  {
	return &DataPack{}
}

//  获取包的头的长度的方法
func (pD *DataPack)GetHandLen() uint32{

	// DataLen uint32(4字节)
	// ID uint32(4字节)
	return 8
}

// 封包方法  dataLen|dataMsgID|data
func (pD *DataPack)Pack(msg ziface.IMessage )([]byte,error){
	// 创建一个存放byte 字节的缓冲
	aDataBuff := bytes.NewBuffer([]byte{})

	// 将 datalen 写入
	if err:=binary.Write(aDataBuff,binary.LittleEndian,msg.GetMsgLen());err != nil {
		return nil, err
	}

	// 将 MsgID 写入
	if err:=binary.Write(aDataBuff,binary.LittleEndian,msg.GetMsgID());err != nil {
		return nil, err
	}

	// 将 data 数据写入
	if err:=binary.Write(aDataBuff,binary.LittleEndian,msg.GetMsgData());err != nil {
		return nil, err
	}

	return  aDataBuff.Bytes(),nil
}

// 拆包方法
// 将包的Head 信息读出来，之后再根据 head 信息里的data 长度，在进行一次读
func (pD *DataPack)Unpack(binaryData []byte)( ziface.IMessage,error){

	// 创建一个从输入二进制数据的 ioReader
	dataBuff :=bytes.NewReader(binaryData)

	// 只解压 head 信息，得到 datalen 和 MsgID
	msg := &Message{}

	// 读datalen
	if err:=binary.Read(dataBuff,binary.LittleEndian,&msg.IDataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err:=binary.Read(dataBuff,binary.LittleEndian,&msg.IID); err != nil {
		return nil, err
	}

	// 判断datalen 是否已经超出了我们运行的最大包长度
	if utils.GlobalObject.IMaxPackageSize >0 &&
			msg.IDataLen > utils.GlobalObject.IMaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg,nil
}










