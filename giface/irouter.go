package giface

//定义路由抽象的方法
type IRouter interface {
	//执行之前
	BeforeHandle(request IRequest)
	//执行主handle
	Handle(request IRequest)
	//执行之后
	AfterHandle(request IRequest)
}
