package gnet

import (
	"github.com/hanjin7278/go-tock/giface"
	"log"
)

type MsgHandler struct {

	//用于存储msgId和router的关系
	Apis map[uint32]giface.IRouter
}

/**
创建MsgHandler
*/
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]giface.IRouter),
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
