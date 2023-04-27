package utils

import (
	"encoding/json"
	"io/ioutil"
	"mmo_game/zinx/ziface"
)

// 存储一切有zinx框架的全局参数，供模块使用
// 一些参数可以通过zinx.json
type Global0bj struct {
	/*
		Server
	*/
	TepServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称
	/*
		Zinx */
	Version          string //当前Zinx的版本号
	MaxConn          int    //当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSiz    uint32 // 业务工作池的work数量
	MaxWorkerPoolSiz uint32 //每个work处理的消息数量
}

/*
定义一个全局的对外G1obalobj
*/
var Global0bject *Global0bj

func init() {
	Global0bject = &Global0bj{
		Name:             "ZinxServer",
		Version:          "v0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		MaxWorkerPoolSiz: 1024,
		WorkerPoolSiz:    10,
	}
	//尝试从zinx.json里加载一些用户自定义从参数
	Global0bject.Reload()
}

// 尝试从zinx.json里加载一些用户自定义从参数
func (g *Global0bj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Global0bject)
	if err != nil {
		panic(err)
	}

}
