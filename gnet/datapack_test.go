package gnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {

	//模拟服务端
	//1 、 创建服务端监听
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("Listen 8888 port err", err)
	}
	go func() {
		// 等待客户端连接
		for {
			accept, err := listen.Accept()
			if err != nil {
				fmt.Println("server accept err", err)
				return
			}
			//有客户端连接进来，则进行拆包操作
			go func(conn net.Conn) {
				//定义
				pack := &DataPack{}
				for {
					//1、读取 head部分内容,
					headBuf := make([]byte, pack.GetHeadLen())
					_, err := io.ReadFull(conn, headBuf)
					if err != nil {
						fmt.Println("server reader header err", err)
						break
					}

					//读取 len 和 id字段
					unpack, err := pack.Unpack(headBuf)

					if err != nil {
						fmt.Println("server unpack err", err)
						break
					}
					//判断是否读取到了内容
					if unpack.GetMessageLen() > 0 {
						//读取data中的内容 返回一个message的指针，使用断言返回Message对象
						msg := unpack.(*Message)
						msg.MsgData = make([]byte, msg.GetMessageLen())
						_, err := io.ReadFull(conn, msg.MsgData)
						if err != nil {
							fmt.Println("server read data err ", err)
						}
						fmt.Printf("reader success : len=%d  id=%d data=%v \n", msg.DataLen, msg.MessageId, string(msg.MsgData))
					}
				}
			}(accept)
		}
	}()

	//模拟客户端

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("客户端连接错误", err)
		return
	}
	dp := NewDataPack()

	msg1 := &Message{MessageId: 1, DataLen: 5, MsgData: []byte{'h', 'e', 'l', 'l', 'o'}}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack err", err)
		return
	}

	msg2 := &Message{MessageId: 2, DataLen: 6, MsgData: []byte{'g', 'o', 'l', 'a', 'n', 'g'}}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack err", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	//阻塞线程
	select {}
}
