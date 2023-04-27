package ziface

/*
连接管理模块抽象层
*/
type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//根据connID获取链接
	Get(connID uint32) (IConnection, error)
	//得到当前连接总数
	Len() int
	//清除并终止所有d连接
	ClearConn()
}
