package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)


// 只是负责测试 datapack 拆包和封包 的单元测试
func TestDataPack( t *testing.T)  {
	/*
	// 模拟服务器
	*/
	// 1 创建socket tcp
	listenner,err := net.Listen("tcp","127.0.0.1:7777")
	if err != nil{
		fmt.Println("server listen err:",err)
		return
	}

	// 创建一个 go 承载 负责从给客户端处理业务
	go func() {
		// 2 从客户端读取数据 菜包处理
		for {
			conn,err := listenner.Accept()
			if err != nil{
				fmt.Println("server Accept err:",err)
				return
			}
			go func( conn net.Conn) {
				// 处理客户端的请求
				// 定义一个拆包的对象 dp
				dp  := NewDataPack()
				for  {
					//  1.第一次从conn 读 把包的 head 读出来
					headData := make([]byte,dp.GetHandLen())
					_,err := io.ReadFull(conn,headData)
					if err != nil{
						fmt.Println("read head error:",err)
						break
					}

					msgHead,err := dp.Unpack(headData)
					if err != nil{
						fmt.Println("dp Unpack error:",err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// 2 第二次从 conn 读，根据 head 中的 datalen 再读取data 内容
						// MSG  是有数据的，需要进行第二次修改
						msg := msgHead.(*Message) // 转换类型
						msg.AData = make([]byte,msg.GetMsgLen())

						// 根据datalen 再次从io流中读取
						_,err :=io.ReadFull(conn,msg.AData)
						if err != nil{
							fmt.Println("io ReadFull error:",err)
							return
						}

						// 完整一个消息已经读取完毕
						fmt.Println("--> Recv MsgID:",msg.IID,"datalen:",msg.IDataLen," Data:" ,msg.AData,"string:",string(msg.AData))
					}
				}


				//
			}(conn)
		}
	}()


	/*
	* 	模拟客户端
	*/

	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if  err != nil{
		fmt.Println("client dial err:",err)
		return
	}

	// 创建一个封包对象 dp
	dp:=NewDataPack()

	// 模拟粘包的过程，封装两个msg一同发送
	msg1:= &Message{
		IID: 1,
		IDataLen: 4,
		AData:[]byte{'z','i','n','x'},
		//AData: []byte("hhhhh"),
	}
	sendData1,err := dp.Pack(msg1)
	if err != nil{
		fmt.Println("send data1 err:",err)
		return
	}
	msg2:= &Message{
		IID: 2,
		IDataLen: 5,
		AData:[]byte{'n','i','h','a','o'},
		//AData: []byte("hhhhh"),
	}
	sendData2,err := dp.Pack(msg2)
	if err != nil{
		fmt.Println("send data2 err:",err)
		return
	}

	// 粘包
	sendData1 = append(sendData1,sendData2...)

	// 一次性发送给服务器
	conn.Write(sendData1)

	// 客户端阻塞
	select {

	}
}












