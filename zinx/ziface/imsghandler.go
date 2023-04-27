package ziface

// 消息管理抽抽象层
type ImsgHandle interface {
	//调度 执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体处理逻辑
	AddRouter(msgOD uint32, router IRouter)
	//消息队列
	StartWorkerPool()
	SendMsgToTakeQueue(request IRequest)
}
