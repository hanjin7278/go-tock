package giface

type IRequest interface {

	//返回当前链接
	GetConnection() IConnection
	//获取数据
	GetData() []byte
	//获取消息Id
	GetMsgId() uint32
}
