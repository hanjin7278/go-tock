package gnet

import (
	"github.com/hanjin7278/go-tock/giface"
)

//BaseRouter 的作用主要是提供默认的调用，客户可继承BaseRouter，复写BeforeHandle,Handle,AfterHandle三个方法
type BaseRouter struct {

}

//执行之前
func (this *BaseRouter) BeforeHandle(request giface.IRequest){}
//执行主handle
func (this *BaseRouter) Handle(request giface.IRequest){}
//执行之后
func (this *BaseRouter) AfterHandle(request giface.IRequest){}


