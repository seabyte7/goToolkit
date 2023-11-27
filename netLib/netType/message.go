package netType

type Message struct {
	Len  uint32 // msg length,include header size
	Data []byte // only data,not include header
}

func NewMessage(msgLen uint32, data []byte) *Message {
	return &Message{
		Len:  msgLen,
		Data: data,
	}
}

func (this *Message) ToMsgBytes() []byte {
	buf := make([]byte, 0, this.Len+MsgHeaderLen)
	buf = append(buf, Uint32ToBytes(this.Len)...)
	buf = append(buf, this.Data...)

	return buf
}

// 判断缓冲区是否是一个完整的数据包
func IsCompleteMessage(data []byte) bool {
	if len(data) < MsgHeaderLen {
		return false
	}

	length := BytesToUint32(data[:MsgHeaderLen])
	return len(data) >= int(length)
}
