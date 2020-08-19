package giface

/**
连接管理的抽象
*/
type IConnManager interface {
	//增加连接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//根据Id查找连接
	GetConnById(connId uint32) (IConnection, error)
	//返回连接总数
	GetConnSize() int
	//清空连接
	Clear()
}
