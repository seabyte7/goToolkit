package tcpLib

import "github.com/seabyte7/goToolkitnetLib/netType"

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
