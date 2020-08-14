package giface

type IMessage interface {
	//获取消息Id
	GetMessageId() uint32
	//获取消息长度
	GetMessageLen() uint32
	//获取消息内容
	GetMsgData() []byte

	//设置消息Id
	SetMessageId(uint32)
	//设置消息的长度
	SetMessageLen(uint32)
	//设置消息内容
	SetMsgData([]byte)
}
