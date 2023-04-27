package ziface

// irequest请求

type IRequest interface {
	//获得当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
	//
	GetMsgID() uint32
}
