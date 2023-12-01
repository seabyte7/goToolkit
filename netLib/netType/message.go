package netType

const (
	CmdHeartbeat uint8 = 1
	CmdMsg       uint8 = 2
)

type ClientServerMsg struct {
	Len  uint32
	Cmd  uint8
	Data []byte
}

func NewClientServerMsg(cmd uint8, data []byte) *ClientServerMsg {
	msgLen := uint32(MsgHeaderLen + len(data))
	return &ClientServerMsg{
		Len:  msgLen,
		Cmd:  cmd,
		Data: data,
	}
}

func (this *ClientServerMsg) ToBytes() []byte {
	buf := make([]byte, 0, MsgHeaderLen+this.Len)
	buf = append(buf, Uint32ToBytes(this.Len)...)
	buf = append(buf, Uint8ToBytes(this.Cmd)...)
	buf = append(buf, this.Data...)

	return buf
}
