package gnet

import (
	"errors"
	"github.com/hanjin7278/go-tock/giface"
	"log"
	"sync"
)

/**
连接管理实现
*/
type ConnManager struct {
	//定义保存连接的集合
	Conns map[uint32]giface.IConnection
	//保护连接操作的读写锁
	Lock sync.RWMutex
}

//创建ConnManager对象方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		Conns: make(map[uint32]giface.IConnection),
	}
}

//增加连接
func (this *ConnManager) Add(conn giface.IConnection) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.Conns[conn.GetConnId()] = conn

	log.Println("[add Connection to Maps connId = ", conn.GetConnId(), " map size is ", len(this.Conns), " Success]")
}

//删除连接
func (this *ConnManager) Remove(conn giface.IConnection) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	delete(this.Conns, conn.GetConnId())

	log.Println("[remove Connection from Maps connId = ", conn.GetConnId(), " map size is ", len(this.Conns), " Success]")
}

//根据Id查找连接
func (this *ConnManager) GetConnById(connId uint32) (giface.IConnection, error) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()

	if conn, ok := this.Conns[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("get connection by connId = " + string(connId) + " not found")
	}
}

//返回连接总数
func (this *ConnManager) GetConnSize() int {
	this.Lock.RLock()
	defer this.Lock.RUnlock()

	log.Println("[get conn size by maps num = ", len(this.Conns), "]")

	return len(this.Conns)
}

//清空连接
func (this *ConnManager) Clear() {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	for k, v := range this.Conns {
		//停止conn工作，并通知读写
		v.Stop()
		//删除连接
		delete(this.Conns, k)
	}
	log.Println("[Clear All conns success]")
}
