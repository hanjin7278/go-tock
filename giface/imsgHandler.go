package giface

type IMsgHandler interface {

	//执行路由方法
	DoMsgHandler(req IRequest)

	//添加路由方法
	AddRouter(msgId uint32, router IRouter)
}
