package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main(){
	fmt.Println("客户端创建连接")
	time.Sleep(1 * time.Second)
	//创建tcp连接客户端
	conn, err := net.Dial("tcp", "0.0.0.0:8888")

	if err != nil {
		log.Fatal("创建tcp客户端连接错误",err)
	}

	for{
		//向服务端发送数据
		conn.Write([]byte("Hello Zink "))

		buf := make([]byte,512)
		_ , err := conn.Read(buf)

		if err != nil {
			log.Fatal("读取内容错误错误",err)
		}
		log.Printf("客户端读取服务端返回内容：%s" , buf)
		time.Sleep(1 * time.Second)
	}
}
