package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/src/zinx/ziface"
)

/*
	连接管理模块实例
*/
type ConnManager struct {
	mapConnections map[uint32] ziface.IConnection	// 管理的连接集合
	connLock	sync.RWMutex // 连接集合的读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		mapConnections:make(map[uint32] ziface.IConnection),
	}
}

//  添加连接
func(pCM *ConnManager)AddConn(conn ziface.IConnection){
	// 加写锁
	pCM.connLock.Lock()
	defer pCM.connLock.Unlock()

	// 将 conn 加入到 connManager 中
	pCM.mapConnections[conn.GetConnID()] = conn
	fmt.Println("connection ID=",conn.GetConnID()," add to ConnManager succ : conn num=",pCM.GetAllConnCount())
}

// 删除连接
func(pCM *ConnManager)RemoveConn(conn ziface.IConnection){
	// 加写锁
	pCM.connLock.Lock()
	defer pCM.connLock.Unlock()

	// 删除连接信息
	delete(pCM.mapConnections,conn.GetConnID())
	fmt.Println("connection ID=",conn.GetConnID()," delete to ConnManager succ : conn num=",pCM.GetAllConnCount())

}

// 根据 ConnID 获取连接
func(pCM *ConnManager)GetConn(iConnID uint32)(	ziface.IConnection,error){
	// 加写锁
	pCM.connLock.RLock()
	defer pCM.connLock.RUnlock()

	if conn,ok := pCM.mapConnections[iConnID];ok {
		// 找到了
		return conn,nil
	}else {
		return nil,errors.New("connection not found !!!")
	}
}

// 得到当前连接的总数
func(pCM *ConnManager)GetAllConnCount() int{
	 return len(pCM.mapConnections)
}

// 清除所有连接
func(pCM *ConnManager)ClearAllConn(){
	// 加写锁
	pCM.connLock.Lock()
	defer pCM.connLock.Unlock()

	// 删除  conn  并停止 工作
	for iConnID,conn := range pCM.mapConnections{
		// 停止
		conn.Stop()

		// 删除
		delete(pCM.mapConnections,iConnID)
	}

	fmt.Println("Clear All Conn Succ num=",pCM.GetAllConnCount())
}

