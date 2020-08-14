package giface

/**
解决TCP粘包的问题，进行拆包封包，序列化TLV 自定义格式
len 长度
code 协议
body 消息体
使用IMessage 封装
head:固定8字节长度获取body内容的长度和消息类别
head:{datalen、messageId} body:{content}
*/

type IDataPack interface {
	//获取Header的长度
	GetHeadLen() uint32
	//封包
	Pack(message IMessage) ([]byte, error)
	//拆包
	Unpack([]byte) (IMessage, error)
}
