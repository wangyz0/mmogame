package znet

import (
	"fmt"
	"mmo_game/zinx/utils"
	"mmo_game/zinx/ziface"
	"strconv"
)

type MsgHandle struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责worker去任务消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作池的work数量
	WorkerPoolSize uint32
}

// 初始化／创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.Global0bject.WorkerPoolSiz, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.Global0bject.MaxWorkerPoolSiz),
	}
}

// 调度／执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID =", request.GetMsgID(), "is NOT FOUND! Need Register!")
		return
	}
	//2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册了
		panic("repeat api ,msgID=" + strconv.Itoa(int(msgID)))
	}
	// 添加msg与API的绑走关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID =", msgID, "succ!")
}

// 启动一个Worker工作池（开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池）
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoo1Size 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		// 1 当前的worker对应的channe1消息队列 开辟空间 第0个worker 就用第0个channe1．．．
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.Global0bject.MaxWorkerPoolSiz)
		// 启动当前的Worker，阻塞等待消息从channe1传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerID, "is started...")
	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//读取到消息 然后执行
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 发送消息到消息队列
func (mh *MsgHandle) SendMsgToTakeQueue(request ziface.IRequest) {
	//   1将消息平均分配到worker队列
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request //将请求发给队列
}
