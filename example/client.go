package main

import (
	"fmt"
	"github.com/hanjin7278/go-tock/gnet"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("客户端创建连接")
	time.Sleep(1 * time.Second)
	//创建tcp连接客户端
	conn, err := net.Dial("tcp", "0.0.0.0:8888")

	if err != nil {
		log.Fatal("创建tcp客户端连接错误", err)
	}

	for {

		dp := gnet.NewDataPack()

		msg := gnet.NewMessage(0, []byte("我是客户端发送的数据"))

		binMsg, err := dp.Pack(msg)
		if err != nil {
			log.Fatal("pack 出现错误", err)
			break
		}

		if _, err := conn.Write(binMsg); err != nil {
			log.Fatal("send data to server err", err)
			break
		}
		//读取服务端返回的数据
		//1、读取Header部分
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			log.Fatal("read head err", err)
			break
		}
		//消息id 和长度
		msgHead, err := dp.Unpack(headData)
		//读取到服务端返回的内容
		if msgHead.GetMessageLen() > 0 {

			data := make([]byte, msgHead.GetMessageLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				log.Fatal("读取data内容出现错误", err)
			}
			msgHead.SetMsgData(data)
			//打印服务端返回的内容
			fmt.Printf("client revc server msgId = %d,len = %d,data = %s\n", msgHead.GetMessageId(), msgHead.GetMessageLen(), string(msgHead.GetMsgData()))
		}
		time.Sleep(1 * time.Second)
	}
}
