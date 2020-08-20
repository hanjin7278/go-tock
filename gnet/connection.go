package gnet

import (
	"errors"
	"fmt"
	"github.com/hanjin7278/go-tock/utils"
	"io"
	"log"
	"net"
	"sync"

	"github.com/hanjin7278/go-tock/giface"
)

/**
定义链接对象内容
*/
type Connection struct {

	//保存当前连接对应的Server对象
	TcpServer giface.IServer
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

	//用于存放连接属性的集合
	ConnProp map[string]interface{}
	//用于保护属性的锁
	PropLock sync.RWMutex
}

/**
初始化连接的方法
*/
func NewConnection(server giface.IServer, conn *net.TCPConn, connId uint32, handle giface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnId:    connId,
		IsClose:   false,
		MsgChan:   make(chan []byte),
		MsgHandle: handle,
		ExitChan:  make(chan bool, 1),
		ConnProp:  make(map[string]interface{}),
	}
	//将当前连接加入到连接管理器中
	server.GetConnMgr().Add(c)

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

	//调用开发者定义的OnConnStart()
	this.TcpServer.CallOnConnStart(this)
}

//停止连接
func (this *Connection) Stop() {
	log.Println("[stop connection exit!!! ] connId=", this.ConnId)
	if this.IsClose {
		return
	}
	this.IsClose = true

	//在连接关闭之前调用开发者定义的OnConnStop()方法
	this.TcpServer.CallOnConnStop(this)

	//关闭Socket连接
	this.Conn.Close()
	//通知Writer关闭
	this.ExitChan <- true
	//关闭管道
	close(this.ExitChan)
	close(this.MsgChan)
	//从连接管理器中移除当前连接
	this.TcpServer.GetConnMgr().Remove(this)

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

//保存属性
func (this *Connection) SetProp(key string, value interface{}) {
	this.PropLock.Lock()
	defer this.PropLock.Unlock()
	this.ConnProp[key] = value
}

//删除属性
func (this *Connection) RemoveProp(key string) {
	this.PropLock.Lock()
	defer this.PropLock.Unlock()
	if _, ok := this.ConnProp[key]; ok {
		delete(this.ConnProp, key)
	} else {
		log.Println("[remove by key ", key, " not in Props]")
	}
}

//获取属性
func (this *Connection) GetProp(key string) (interface{}, error) {
	this.PropLock.RLock()
	defer this.PropLock.RUnlock()
	if _, ok := this.ConnProp[key]; ok {
		return this.ConnProp[key], nil
	} else {
		log.Println("[remove by key ", key, " not in Props]")
		return nil, errors.New("[get by key not in Props]")
	}
}

//清空属性
func (this *Connection) ClearProp() {
	this.PropLock.Lock()
	defer this.PropLock.Unlock()

	for k, _ := range this.ConnProp {
		delete(this.ConnProp, k)
	}
	log.Println("[Clear All Properties Success]")
}
