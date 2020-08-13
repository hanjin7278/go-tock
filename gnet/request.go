package gnet

import (
	"github.com/hanjin7278/go-tock/giface"
)

//封装Request请求
type Request struct {
	//当前链接
	conn giface.IConnection
	//请求的数据
	data []byte

}

//返回当前链接
func (this *Request) GetConnection() giface.IConnection{
	return this.conn
}
//获取数据
func (this *Request) GetData() []byte{
	return this.data
}
