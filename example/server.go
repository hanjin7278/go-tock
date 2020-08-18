package main

import (
	"fmt"
	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/gnet"
	"log"
)

type MyRouter struct {
	gnet.BaseRouter
}

//执行主handle
func (this *MyRouter) Handle(request giface.IRequest) {
	fmt.Printf("Server revc client msgId = %d,data = %s\n", request.GetMsgId(), string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping....ping...ping"))
	if err != nil {
		log.Fatal("send data to client err", err)
	}
}

type TowRouter struct {
	gnet.BaseRouter
}

//执行主handle
func (this *TowRouter) Handle(request giface.IRequest) {
	fmt.Printf("Server revc client msgId = %d,data = %s\n", request.GetMsgId(), string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Welcome Tow Handler go-tock"))
	if err != nil {
		log.Fatal("send data to client err", err)
	}
}

func main() {
	server := gnet.NewServer()
	//添加多个自定义路由
	server.AddRouter(0, &MyRouter{})
	server.AddRouter(1, &TowRouter{})
	server.Run()
}
