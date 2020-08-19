package giface

//定义抽象的接口方法
type IServer interface {
	//停止服务器
	Stop()

	//运行服务器
	Run()

	//增加Router
	AddRouter(msgId uint32, router IRouter)

	//获取连接管理器
	GetConnMgr() IConnManager

	//注册OnConnStart()方法
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop()方法
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStart()方法
	CallOnConnStart(conn IConnection)
	//调用OnConnStop()方法
	CallOnConnStop(conn IConnection)
}
