package tcpLib

import "goToolkit/netLib/netType"

type ClientServerMsg struct {
	Session *TcpSession
	netType.NetMsg
}

func NewClientServerMsg(sessionPtr *TcpSession, msg *netType.NetMsg) *ClientServerMsg {
	return &ClientServerMsg{
		Session: sessionPtr,
		NetMsg:  *msg,
	}
}
