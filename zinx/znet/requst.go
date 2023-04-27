package znet

import "mmo_game/zinx/ziface"

type Request struct {
	//已经和客户端获得当前连接
	conn ziface.IConnection
	// 得到请求的消息数据
	msg ziface.IMessage
}

//获得当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
