package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/hanjin7278/go-tock/giface"
	"github.com/hanjin7278/go-tock/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取Header的长度 header固定长度：len (uint32) 4字节 + id (uint32)4字节 = 8 字节
func (this *DataPack) GetHeadLen() uint32 {
	return 8
}

/**
封包
|len|id|data| 这种格式写入
*/
func (this *DataPack) Pack(message giface.IMessage) ([]byte, error) {
	//创建带缓冲的buffer
	buffer := bytes.NewBuffer([]byte{})
	//写入Len长度
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMessageLen()); err != nil {
		return nil, err
	}
	//写入id
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMessageId()); err != nil {
		return nil, err
	}
	//写入data
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

//拆包
func (this *DataPack) Unpack(data []byte) (giface.IMessage, error) {

	//创建一个二进制流读取
	reader := bytes.NewReader(data)

	//定义Message对象
	msg := &Message{}
	//读取长度
	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取id
	if err := binary.Read(reader, binary.LittleEndian, &msg.MessageId); err != nil {
		return nil, err
	}

	//判断包大小
	if utils.GlobalConfigObj.MaxPackageSize > 0 && msg.DataLen > utils.GlobalConfigObj.MaxPackageSize {
		return nil, errors.New("传输的包过大！")
	}
	return msg, nil
}
