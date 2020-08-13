package giface

//定义抽象的接口方法
type IServer interface {

	//启动服务器
	Start()

	//停止服务器
	Stop()

	//运行服务器
	Run()

	//增加Router
	AddRouter(router IRouter)
}
