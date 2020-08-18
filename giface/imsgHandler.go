package giface

type IMsgHandler interface {

	//执行路由方法
	DoMsgHandler(req IRequest)

	//添加路由方法
	AddRouter(msgId uint32, router IRouter)
	//启动工作池
	StartWorkerPool()
	//发送消息到队列
	SendMsgToTaskQueue(req IRequest)
}
