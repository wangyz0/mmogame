package znet

import (
	"errors"
	"fmt"
	"io"
	"mmo_game/zinx/utils"
	"mmo_game/zinx/ziface"
	"net"
	"sync"
)

// 链接模块

type Connection struct {
	// 当前的conn属于哪个server
	TcpServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的Id
	ConnID uint32
	// 当前链接的状态
	isClosed bool
	// 告知当前连接已经停止的channel
	ExitChan chan bool
	// 无缓存管道  用于读写goroutine之间通信
	msgChan chan []byte
	// 消息的管理msgID和对应的处理业务
	MsgHandler ziface.ImsgHandle
	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.ImsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)

	return c
}
func (c *Connection) StartRead() {
	fmt.Println("开始读消息connID=", c.ConnID)
	defer c.Stop()
	defer fmt.Println("reader 退出")
	for {
		// buf := make([]byte, 512)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Printf("err: %v\n", err)
		// 	continue
		// }
		// // // 调用当前
		// // if err := c.handleAPI(c.Conn, buf[:cnt], cnt); err != nil {
		// // 	fmt.Printf("err: %v\n", err)
		// // 	break
		// // }
		// // 得到当前conn数据的request请求数据
		// 创建拆包对象
		dp := NewDataPack()
		// 读取客户端的msg head 8字节
		headData := make([]byte, dp.GetHeadLen()) //8个字节
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		// 拆包  得到msgID 和msgDataLEn 放入masg消息
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		// 根据datalen读取data 放在msg.data里
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			// fmt.Println("读到消息", string(data))
			if err != nil {
				fmt.Printf("err: %v\n", err)
				break
			}
			msg.SetData(data)
		}

		req := Request{
			conn: c,
			msg:  msg,
		}
		// 判断是否开启工作池子
		if utils.Global0bject.WorkerPoolSiz > 0 {
			c.MsgHandler.SendMsgToTakeQueue(&req)
		} else {
			// 从里有中 找到注册绑定的router调用
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// 写消息的，专门发送给客户端消息
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn writer exit")
	for {
		// 阻塞等待消息

		select {
		case data := <-c.msgChan:
			// fmt.Println("msgChan管道收到data", string(data))
			//将data写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("send data err: %v\n", err)
			}
		case b := <-c.ExitChan: //如果可读  说明连接要关闭了
			fmt.Printf("b: %v\n", b)
			return
		}
	}
}

// 开始连接
func (c *Connection) Start() {
	fmt.Println("conn start() conn.id=", c.ConnID)

	//   启动写数据业务
	fmt.Println("开始写")
	go c.StartWriter()
	// 按照开发者所传来的方法调用hook函数
	c.TcpServer.Cal10nConnStart(c)
	//    启动当前连接读数据
	fmt.Println("start函数")
	c.StartRead()

}

// 停止链接
func (c *Connection) Stop() {
	fmt.Println("conn stop() conn.id=", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 按照开发者所传来的方法调用hook函数
	c.TcpServer.Cal10nConnStop(c)
	// 关闭连接
	c.Conn.Close()
	//告知其他grotine 连接要关闭
	c.ExitChan <- true
	// 关闭管道
	close(c.ExitChan)
	close(c.msgChan)
	//将链接从连接池删除
	c.TcpServer.GetConnMgr().Remove(c)
}

// 获取链接绑定的链接ID
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取客户端的tcp状态  ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// // 发送数据给客户端
//
//	func (c *Connection) Send(data []byte) error {
//		return nil
//	}
//
// 提供一个sendmsg方法  将发送给客户端的数据先封包
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	// fmt.Println("Send收到的data:", string(data))
	if c.isClosed == true {
		return errors.New("连接关闭")
	}
	// 封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	// fmt.Println("发送的binaryMsg", string(binaryMsg))
	c.msgChan <- binaryMsg
	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	// 添加一个属性
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("没有该连接属性")
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
