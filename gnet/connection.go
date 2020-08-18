package gnet

import (
	"errors"
	"fmt"
	"github.com/hanjin7278/go-tock/utils"
	"io"
	"log"
	"net"

	"github.com/hanjin7278/go-tock/giface"
)

/**
定义链接对象内容
*/
type Connection struct {
	//当前的链接
	Conn *net.TCPConn
	//链接id
	ConnId uint32
	//当前链接是否关闭
	IsClose bool
	//是否退出
	ExitChan chan bool

	//用于reader、writer 之间通信的管道
	MsgChan chan []byte

	//增加路由成员
	MsgHandle giface.IMsgHandler
}

/**
初始化连接的方法
*/
func NewConnection(conn *net.TCPConn, connId uint32, handle giface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnId:    connId,
		IsClose:   false,
		MsgChan:   make(chan []byte),
		MsgHandle: handle,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

/**
启动读的协程
*/
func (this *Connection) StartReader() {
	log.Println("[Connection:start read data running] connId = ", this.ConnId)
	defer log.Println("[Reader close connection,Exit!!!]Conn=", this.Conn.RemoteAddr().String())
	defer this.Stop()
	for {
		dp := NewDataPack()
		headBuf := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(this.Conn, headBuf)
		if err != nil {
			fmt.Println("server reader header err", err)
			break
		}
		//读取 len 和 id字段
		msg, err := dp.Unpack(headBuf)
		if err != nil {
			fmt.Println("server msg err", err)
			break
		}
		//判断是否读取到了内容
		var data []byte
		if msg.GetMessageLen() > 0 {
			//读取data中的内容 返回一个message的指针，使用断言返回Message对象
			data = make([]byte, msg.GetMessageLen())
			_, err := io.ReadFull(this.Conn, data)
			if err != nil {
				fmt.Println("server read data err ", err)
				break
			}
			//将数据设置到Message对象的data属性中
			msg.SetMsgData(data)
			r := Request{
				conn: this,
				msg:  msg,
			}
			//调用用户自定义的Handle
			if utils.GlobalConfigObj.WorkerPoolSize > 0 {
				//已经开启工作池，交给工作池处理
				this.MsgHandle.SendMsgToTaskQueue(&r)
			} else {
				//还用原始的处理
				go this.MsgHandle.DoMsgHandler(&r)
			}

		}
	}
}

/**
用于回写到客户端的goroutine
*/
func (this *Connection) StartWriter() {
	log.Println("[Connection:start writer data running] connId = ", this.ConnId)
	defer log.Println("[Writer closed connection,Exit!!!] conn=", this.Conn.RemoteAddr().String())

	//不断监听管道里的内容
	for {
		select {
		case data := <-this.MsgChan:
			//从管道中读取到了数据
			if _, err := this.Conn.Write(data); err != nil {
				log.Fatal("writer to client err", err)
			}
		case <-this.ExitChan:
			//读取到了退出消息
			return
		}
	}
}

//启动连接
func (this *Connection) Start() {
	go this.StartReader()
	go this.StartWriter()
}

//停止连接
func (this *Connection) Stop() {
	log.Println("[stop connection exit!!! ] connId=", this.ConnId)
	if this.IsClose {
		return
	}
	this.IsClose = true
	//关闭Socket连接
	this.Conn.Close()
	//通知Writer关闭
	this.ExitChan <- true
	//关闭管道
	close(this.ExitChan)
	close(this.MsgChan)
	log.Println("ConnId = ", this.ConnId, " closed")
}

//获取当前连接绑定的Socket
func (this *Connection) GetSocketConn() *net.TCPConn {
	return this.Conn
}

//获取当前连接的Id
func (this *Connection) GetConnId() uint32 {
	return this.ConnId
}

//获取远程客户端的ip和端口
func (this *Connection) RemoteAddr() net.Addr {
	return this.Conn.RemoteAddr()
}

//发送数据
func (this *Connection) SendMsg(msgId uint32, data []byte) error {
	if this.IsClose == true {
		log.Fatal("Connection is Closed ")
		return errors.New("Connection is Closed")
	}
	dp := NewDataPack()
	binMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		log.Fatal("Server Pack error", err)
	}
	//将数据发送到写的管道
	this.MsgChan <- binMsg

	return nil
}
