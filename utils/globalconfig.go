package utils

import (
	"encoding/json"
	"github.com/hanjin7278/go-tock/giface"
	"io/ioutil"
)

type GlobalConfig struct {
	TcpServer giface.IServer //构建的Server对象

	ServerName string //服务名称

	Host string //绑定的ip

	Port int //绑定的端口号

	MaxConn int //允许客户端连接的最大数

	MaxPackageSize uint32 //客户端一次请求包的最大值

	Version string //go-tock 版本号
}

var GlobalConfigObj *GlobalConfig

/**
加载配置文件
*/
func (g *GlobalConfig) Reload() {
	data, err := ioutil.ReadFile("./config/go-tock.json")
	PanicError(err)
	err = json.Unmarshal(data, &GlobalConfigObj)
	PanicError(err)
}

/**
初始化全局配置对象
*/
func init() {
	GlobalConfigObj = &GlobalConfig{
		ServerName:     "go-tock-server",
		Host:           "0.0.0.0",
		Port:           8888,
		MaxConn:        2000,
		MaxPackageSize: 4096,
		Version:        "V0.1",
	}
	//刷新配置文件信息
	//GlobalConfigObj.Reload()
}
