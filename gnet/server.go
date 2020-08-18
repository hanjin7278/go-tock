package gnet

import (
	"fmt"
	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/utils"
	"log"
	"net"
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
}

//实现IServer的Start方法
func (this *Server) Start() {
	go func() {
		log.Printf("[Start] Server Listenner at IP: %s Port:%d", this.IP, this.Port)

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
			//调用Connection模块读取数据
			c := NewConnection(acceptTCP, cid, this.MsgHandle)
			cid++
			go c.Start()
		}
	}()
}

//实现IServer的Stop方法
func (this *Server) Stop() {

}

//实现IServer的Run方法
func (this *Server) Run() {
	this.Start()
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
	}
	return s
}
