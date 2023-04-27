package ziface

import "net"

type IConnection interface {
	// 启动链接
	Start()
	//  停止链接
	Stop()
	// 获取链接绑定的链接ID
	GetTCPConnection() *net.TCPConn
	// 获取当前链接模块的链接ID
	GetConnID() uint32
	// 获取客户端的tcp状态  ip port
	RemoteAddr() net.Addr
	// 发送数据给客户端
	SendMsg(msgId uint32, data []byte) error
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

// 定义一个处理链接的业务方法
type HandleFunc func(*net.TCPConn, []byte, int) error
