package main

import (
	"fmt"
	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/gnet"
)

type MyRouter struct {
	gnet.BaseRouter
}

//执行之前
func (this *MyRouter) BeforeHandle(request giface.IRequest){
	fmt.Println("Server 执行之前 BeforeHandle")
	request.GetConnection().Send([]byte("Before Ping ....\n"))
}
//执行主handle
func (this *MyRouter) Handle(request giface.IRequest){
	fmt.Println("Server 执行 Handle")
	request.GetConnection().Send([]byte("Handle Ping ....\n"))
}
//执行之后
func (this *MyRouter) AfterHandle(request giface.IRequest){
	fmt.Println("Server 执行之后 AfterHandle")
	request.GetConnection().Send([]byte("AfterHandle Ping ....\n"))
}


func main(){
	server := gnet.NewServer("123")
	r := MyRouter{}
	server.AddRouter(&r)
	server.Run()
}
