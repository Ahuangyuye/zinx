package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/src/zinx/znet"
)

/*
模拟客户端
*/

func main(){
	fmt.Println("client start ...")
	time.Sleep(1* time.Second)
	// 1 直接连接远程服务器 得到 conn
	pConn,err := net.Dial("tcp","127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err:",err)
		return
	}


	// 2 连接调用 write 写数据
	for {
		//_,err := pConn.Write([]byte("hello zinx  v0.5..."))
		//if err != nil {
		//	fmt.Println("client start err:",err)
		//	return
		//}
		//
		//aBuf := make([]byte,512)
		//iCount ,err := pConn.Read(aBuf)
		//if err != nil {
		//	fmt.Println("client start err:",err)
		//	return
		//}
		//
		//fmt.Printf("server call back:%sback,iCount=%d\n",aBuf,iCount)


		// 发送封包的message 消息
		objDP := znet.NewDataPack()
		binarryMsg,err := objDP.Pack(znet.NewMsgPackage(0,[]byte("hello zinx  v0.5...")))
		if err != nil{
			fmt.Println("objDP Pack err:",err)
			return
		}

		if  _,err:=pConn.Write(binarryMsg);err != nil{
			fmt.Println("pConn Write err:",err)
			return
		}

		// 收数据

		binarryHead := make([]byte,objDP.GetHandLen())
		if _,err :=  io.ReadFull(pConn,binarryHead);err != nil{
			fmt.Println("pConn io ReadFull err:",err)
			break
		}
		msgHead,err := objDP.Unpack(binarryHead)
		if err != nil{
			fmt.Println("objDP  Unpack err:",err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// 第二次读取具体数据
			msg:= msgHead.(*znet.Message)
			msg.AData = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(pConn,msg.AData);err != nil{
				fmt.Println("pConn io ReadFull err:",err)
				break
			}
			fmt.Println("--> recv server msg: id=",msg.IID,
				" len=",msg.IDataLen,
				" data=",string(msg.AData))
		}

		time.Sleep(1* time.Second)
	}

}



