package znet

type Message struct {
	ID      uint32
	DataLen uint32
	Data    []byte
}

// 创建一个message
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.ID
}

// 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.ID = id
}

// 置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 设置消息的长度
func (m *Message) SetDataLen(datalen uint32) {
	m.DataLen = datalen
}
