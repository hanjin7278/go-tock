package gnet

import (
	"errors"
	"fmt"
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
	//增加路由成员
	Router giface.IRouter
}

/**
初始化连接的方法
*/
func NewConnection(conn *net.TCPConn, connId uint32, router giface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnId:   connId,
		IsClose:  false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

func (this *Connection) StartReader() {
	log.Println("ConnId = ", this.ConnId, "Connection:开始读取数据")
	defer log.Println("ConnId=", this.ConnId, " 正在关闭连接")
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
			go func(req giface.IRequest) {
				this.Router.BeforeHandle(req)
				this.Router.Handle(req)
				this.Router.AfterHandle(req)
			}(&r)
		}
	}
}

//启动连接
func (this *Connection) Start() {
	go this.StartReader()
}

//停止连接
func (this *Connection) Stop() {
	log.Println("正在关闭连接 ConnId = ", this.ConnId)
	if this.IsClose {
		return
	}
	this.IsClose = true
	//关闭Socket连接
	this.Conn.Close()
	//关闭管道
	close(this.ExitChan)
	log.Println("连接ConnId = ", this.ConnId, "已经关闭")
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
	//将数据发送到客户端
	if _, err := this.Conn.Write(binMsg); err != nil {
		log.Printf("发送客户端数据出现错误", err)
		return errors.New("发送客户端数据出现错误")
	}
	return nil
}
