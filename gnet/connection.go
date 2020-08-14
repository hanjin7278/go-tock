package gnet

import (
	"errors"
	"log"
	"net"

	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/utils"
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
	log.Println("ConnId = ", this.ConnId, " 开始读取数据")
	defer log.Println("ConnId=", this.ConnId, " 正在关闭连接")
	defer this.Stop()
	for {
		buf := make([]byte, utils.GlobalConfigObj.MaxPackageSize)
		_, err := this.Conn.Read(buf)
		if err != nil {
			log.Fatal("读取错误 ConnId = ", this.ConnId, err)
			continue
		}

		r := Request{
			conn: this,
			data: buf,
		}
		//调用用户自定义的Handle
		go func(req giface.IRequest) {
			this.Router.BeforeHandle(req)
			this.Router.Handle(req)
			this.Router.AfterHandle(req)
		}(&r)
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
func (this *Connection) Send(data []byte) error {
	if _, err := this.Conn.Write(data); err != nil {
		log.Printf("发送客户端数据出现错误", err)
		return errors.New("发送客户端数据出现错误")
	}
	return nil
}
