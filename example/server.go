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
	//消息Id是0的
	if request.GetMsgId() == 0 {
		err := request.GetConnection().SendMsg(1, []byte("ping....ping...ping"))
		if err != nil {
			log.Fatal("send data to client err", err)
		}
	}

}

func main() {
	server := gnet.NewServer()
	r := MyRouter{}
	server.AddRouter(0, &r)
	server.Run()
}
