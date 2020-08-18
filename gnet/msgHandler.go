package gnet

import (
	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/utils"
	"log"
)

type MsgHandler struct {

	//用于存储msgId和router的关系
	Apis map[uint32]giface.IRouter

	//增加消息队列
	TaskQueue []chan giface.IRequest

	//增加worker工作池数量
	WorkerPoolSize uint32
}

/**
创建MsgHandler
*/
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]giface.IRouter),
		WorkerPoolSize: utils.GlobalConfigObj.WorkerPoolSize,
		TaskQueue:      make([]chan giface.IRequest, utils.GlobalConfigObj.WorkerPoolSize),
	}
}

//执行路由方法
func (this *MsgHandler) DoMsgHandler(req giface.IRequest) {
	router, ok := this.Apis[req.GetMsgId()]
	if ok {
		router.BeforeHandle(req)
		router.Handle(req)
		router.AfterHandle(req)
	} else {
		panic("消息id = [" + string(req.GetMsgId()) + "] 对应的路由不存在")
	}
}

//添加路由方法
func (this *MsgHandler) AddRouter(msgId uint32, router giface.IRouter) {

	if _, ok := this.Apis[msgId]; ok {
		log.Fatal("路由方法已经注册")
		return
	}
	this.Apis[msgId] = router
}

/**
启动工作池(只执行一次)
*/
func (this *MsgHandler) StartWorkerPool() {

	for i := 0; i < int(this.WorkerPoolSize); i++ {
		this.TaskQueue[i] = make(chan giface.IRequest, utils.GlobalConfigObj.MaxWorkerTaskLen)
		go this.startOnWorker(i, this.TaskQueue[i])
	}

}

/**
启动任务
*/
func (this *MsgHandler) startOnWorker(workerId int, task chan giface.IRequest) {
	log.Printf("[ startOnWorker workerId=%d ]\n", workerId)
	//阻塞等待消息到来
	for {
		select {
		case t := <-task:
			this.DoMsgHandler(t)
		}
	}
}

/**
发送请求到消息队列中
其中使用平均取余算法
*/
func (this *MsgHandler) SendMsgToTaskQueue(request giface.IRequest) {
	workerId := request.GetConnection().GetConnId() % this.WorkerPoolSize
	log.Printf("[add task connId = %d , msgId=%d , to workerId=%d ]\n", request.GetConnection().GetConnId(),
		request.GetMsgId(), workerId)
	this.TaskQueue[int(workerId)] <- request
}
