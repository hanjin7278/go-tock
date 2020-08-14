package gnet

type Message struct {
	MessageId uint32 //消息Id

	DataLen uint32 //消息的长度

	MsgData []byte //消息内容
}

//获取消息Id
func (this *Message) GetMessageId() uint32 {
	return this.MessageId
}

//获取消息长度
func (this *Message) GetMessageLen() uint32 {
	return this.DataLen
}

//获取消息内容
func (this *Message) GetMsgData() []byte {
	return this.MsgData
}

//设置消息Id
func (this *Message) SetMessageId(messageId uint32) {
	this.MessageId = messageId
}

//设置消息的长度
func (this *Message) SetMessageLen(dataLen uint32) {
	this.DataLen = dataLen
}

//设置消息内容
func (this *Message) SetMsgData(data []byte) {
	this.MsgData = data
}
