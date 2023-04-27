package ziface

// 定义一个服务器接口
type IServer interface {
	//1启动服务器
	Start()
	// 2停止服务器
	Stop()
	// 3 运行服务器
	Server()
	// 路由功能  给当前服务注册一个路由方法  供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	GetConnMgr() IConnManager
	//
	//注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(conneciton IConnection))
	//注册OnConnStop钩子函数的方法
	SetOnConnStop(func(conneciton IConnection))
	//调用OnConnStart钩子函数的方法
	Cal10nConnStart(conneciton IConnection)
	//调用OnConnStop钩子函数的方法
	Cal10nConnStop(conneciton IConnection)
}
