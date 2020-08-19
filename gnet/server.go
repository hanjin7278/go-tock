package gnet

import (
	"fmt"
	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/utils"
	"log"
	"net"
	"time"
)

//定义IServer接口的实现
type Server struct {
	//服务器名称
	Name string
	//绑定IP的版本
	IPVersion string
	//绑定的IP地址
	IP string
	//监听的端口
	Port int
	//增加路由成员
	MsgHandle giface.IMsgHandler
	//增加连接管理器
	ConnMgr giface.IConnManager
	//连接创建后调用的方法
	OnConnStart func(giface.IConnection)
	//连接销毁之前调用的方法
	OnConnStop func(giface.IConnection)
}

//实现IServer的Start方法
func (this *Server) start() {
	go func() {
		this.MsgHandle.StartWorkerPool()
		time.Sleep(10 * time.Millisecond)
		log.Println("[Start WorkerPool Success]")

		log.Printf("[Start] Server Listenner at IP: %s Port:%d MaxConnection: %d\n", this.IP, this.Port, utils.GlobalConfigObj.MaxConn)

		addr, err := net.ResolveTCPAddr(this.IPVersion, fmt.Sprintf("%s:%d", this.IP, this.Port))
		if err != nil {
			log.Fatal("Server Start err", err)
			return
		}
		//启动监听tcp
		tcp, err := net.ListenTCP(this.IPVersion, addr)
		if err != nil {
			log.Fatal("listenner Tcp err", err)
			return
		}
		log.Printf("Start go-tock is Success , Server Name is %s", this.Name)

		var cid uint32 = 0
		//等待客户端连接
		for {
			acceptTCP, err := tcp.AcceptTCP()
			if err != nil {
				log.Fatal("accept err", err)
				continue
			}

			//创建连接的时候判断当前连接数是否大于用于配置的数量
			if this.ConnMgr.GetConnSize() >= utils.GlobalConfigObj.MaxConn {
				//TODO 以后加入提示关闭连接的错误信息返回
				log.Println("[too many Connection max connection num is ", utils.GlobalConfigObj.MaxConn, "]")
				acceptTCP.Close()
				continue
			}

			//调用Connection模块读取数据
			c := NewConnection(this, acceptTCP, cid, this.MsgHandle)
			cid++
			go c.Start()
		}
	}()
}

//获取当前连接管理器
func (this *Server) GetConnMgr() giface.IConnManager {
	return this.ConnMgr
}

//实现IServer的Stop方法
func (this *Server) Stop() {
	//清理连接，回收资源
	this.ConnMgr.Clear()
	log.Println("[Server ", this.Name, " Stop Success!]")
}

//实现IServer的Run方法
func (this *Server) Run() {
	this.start()
	//TODO 以后可以处理其它的业务

	//阻塞
	select {}
}

//添加router设置
func (this *Server) AddRouter(msgId uint32, router giface.IRouter) {
	this.MsgHandle.AddRouter(msgId, router)
}

/**
创建Server，返回Server的实例
*/
func NewServer() *Server {
	s := &Server{
		Name:      utils.GlobalConfigObj.ServerName,
		IPVersion: "tcp4",
		IP:        utils.GlobalConfigObj.Host,
		Port:      utils.GlobalConfigObj.Port,
		MsgHandle: NewMsgHandler(),
		ConnMgr:   NewConnManager(),
	}
	return s
}

//注册OnConnStart()方法
func (this *Server) SetOnConnStart(hookFun func(conn giface.IConnection)) {
	this.OnConnStart = hookFun
}

//注册OnConnStop()方法
func (this *Server) SetOnConnStop(hookFun func(conn giface.IConnection)) {
	this.OnConnStop = hookFun
}

//调用OnConnStart()方法
func (this *Server) CallOnConnStart(conn giface.IConnection) {
	if this.OnConnStart != nil {
		log.Println("====>[CallOnConnStart() .....]")
		this.OnConnStart(conn)
	}
}

//调用OnConnStop()方法
func (this *Server) CallOnConnStop(conn giface.IConnection) {
	if this.OnConnStop != nil {
		log.Println("====>[CallOnConnStop() .....]")
		this.OnConnStop(conn)
	}
}
