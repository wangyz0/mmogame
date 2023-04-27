package znet

import (
	"fmt"
	"mmo_game/zinx/utils"
	"mmo_game/zinx/ziface"
	"net"
	"time"
	//"zinx-mmo/zinx/ziface"
)

// IServer 接口实现的结构体
type Server struct {
	//1服务器名称
	Name string
	//2服务器ip版本
	IPVersion string
	//3服务器监听ip
	IP string
	//4服务器监听端口
	Port int
	// 当前的Server的消息模块  绑定magid处理业务
	MsgHandler ziface.ImsgHandle
	//连接管理器
	ConnMgr ziface.IConnManager
	// 链接启动和销毁后执行的函数
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

// // 定义当前的handlerapi  这个将来有使用框架的用户来自定
// func CallBackClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	fmt.Println("回写客户端")
// 	_, err := conn.Write(data[:cnt])
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return errors.New("回写客户端错误")
// 	}
// 	return nil

// }

// 1启动服务器
func (s *Server) Start() {
	//开启消息队列
	s.MsgHandler.StartWorkerPool()
	go func() {
		fmt.Printf("开始监听 ip:=%s,port=%d\n", s.IP, s.Port)
		// 获取一个TCP的ADDR

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("net.ResolveTCPAddr err: %v\n", err)
			return
		}
		// 2监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("net.ListenTCP err: %v\n", err)
			return
		}
		defer listenner.Close()

		fmt.Printf("开始监听%s\n", s.Name)
		// 3阻塞的等待客户端连接  处理客户端业务
		var cid uint32 = 0
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("listen.accept连接失败", err)
				continue
			}
			fmt.Println("得到连接")
			// 如果超过最大连接数量 关闭链接
			if s.ConnMgr.Len() > utils.Global0bject.MaxConn {
				conn.Close()
				continue
			}
			// 将处理链接的业务方法和conn绑定得到链接模块
			cid++

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			fmt.Println("启动当前连接业务处理")
			// 启动当前连接业务处理
			go dealConn.Start()
		}
	}()
}

// 2停止服务器
func (s *Server) Stop() {
	// TODO  释放一些资源 停止
	//删除链接
	s.ConnMgr.ClearConn()
}

// 3 运行服务器
func (s *Server) Server() {
	// 启动服务
	s.Start()
	// TODO 做一些服务器启动后的额外业务
	// 阻塞  因为start里是用协程跑的
	time.Sleep(time.Hour)

}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("add router 成功")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 初始化server
func NewServer(name string) *Server {
	s := &Server{
		Name:       utils.Global0bject.Name,
		IPVersion:  "tcp4",
		IP:         utils.Global0bject.Host,
		Port:       utils.Global0bject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookfunk func(ziface.IConnection)) {
	s.OnConnStart = hookfunk
}

// 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookfunc func(conneciton ziface.IConnection)) {
	s.OnConnStop = hookfunc
}

// 调用OnConnStart钩子函数的方法
func (s *Server) Cal10nConnStart(conneciton ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("调用OnConnStart")
		s.OnConnStart(conneciton)
	}
}

// 调用OnConnStop钩子函数的方法
func (s *Server) Cal10nConnStop(conneciton ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("调用OnConnStop")
		s.OnConnStop(conneciton)
	}
}
