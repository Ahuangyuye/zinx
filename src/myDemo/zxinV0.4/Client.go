package main

import (
	"fmt"
	"net"
	"time"
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
		_,err := pConn.Write([]byte("hello zinx  v0.2..."))
		if err != nil {
			fmt.Println("client start err:",err)
			return
		}

		aBuf := make([]byte,512)
		iCount ,err := pConn.Read(aBuf)
		if err != nil {
			fmt.Println("client start err:",err)
			return
		}

		fmt.Printf("server call back:%sback,iCount=%d\n",aBuf,iCount)

		time.Sleep(1* time.Second)
	}

}



