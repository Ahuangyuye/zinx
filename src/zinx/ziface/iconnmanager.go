package ziface

/*
	连接管理模块抽象层
*/

type IConnManager interface {
	//  添加连接
	AddConn(conn IConnection)

	// 删除连接
	RemoveConn(conn IConnection)

	// 根据 ConnID 获取连接
	GetConn(iConnID uint32)(IConnection,error)

	// 得到当前连接的总数
	GetAllConnCount() int

	// 清除所有连接
	ClearAllConn()
}
