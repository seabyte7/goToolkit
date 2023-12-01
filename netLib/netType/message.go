package netType

/*
 * 1. 通信协议
服务器接受tcp链接
1:客户端链接成功后,主动发起登录请求,超过10秒未登录,服务器主动断开链接
未登录链接会将数据包发送给login服务器解析(https)，解析成功后返回登录数据包给客户端
2:客户端登录成功后,服务器主动发起心跳包,客户端回复心跳包,超过10秒未回复,服务器主动断开链接
 * 2. 通信数据结构

*/

const (
	CmdHeartbeat = 1
)

type ServerHeader struct {
	SrcSessionID uint32
	Cmd          int8
	CmdValue     int8
}

type ServerMessage struct {
	Header ServerHeader
	Body   *MessageBody
}

func NewServerMessage(srcSessionID uint32, cmd, cmdValue int8, data []byte) *ServerMessage {
	bodyPtr := NewMessageBody(uint32(len(data))+MsgHeaderLen, data)

	return &ServerMessage{
		Header: ServerHeader{
			SrcSessionID: srcSessionID,
			Cmd:          cmd,
			CmdValue:     cmdValue,
		},
		Body: bodyPtr,
	}
}

type MessageBody struct {
	SrcSessionID  uint64
	DestSessionID uint64
	Len           uint32 // msg length,include header size
	Data          []byte // only data,not include header
}

func NewMessageBody(msgLen uint32, data []byte) *MessageBody {
	return &MessageBody{
		Len:  msgLen,
		Data: data,
	}
}

func (this *MessageBody) ToBytes() []byte {
	buf := make([]byte, 0, this.Len+MsgHeaderLen)
	buf = append(buf, Uint32ToBytes(this.Len)...)
	buf = append(buf, this.Data...)

	return buf
}

type Client2Server struct {
}
