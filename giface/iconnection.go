package giface

import "net"

type IConnection interface {

	//启动连接
	Start()

	//停止连接
	Stop()

	//获取当前连接绑定的Socket
	GetSocketConn() *net.TCPConn

	//获取当前连接的Id
	GetConnId() uint32

	//获取远程客户端的ip和端口
	RemoteAddr() net.Addr

	//发送数据
	SendMsg(msgId uint32, data []byte) error

	//保存属性
	SetProp(key string, value interface{})
	//删除属性
	RemoveProp(key string)
	//获取属性
	GetProp(key string) (interface{}, error)
	//清空属性
	ClearProp()
}

//定义处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
