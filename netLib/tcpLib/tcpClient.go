package tcpLib

import (
	"goToolkit/netLib/netType"
	. "goToolkit/protocol"
	"net"
)

type TcpClient struct {
	Name          string
	OutMsgChannel chan *netType.Message

	session *TcpSession
}

func DialTcpServer(name, addr string) (*TcpClient, Result) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	msgChannel := make(chan *netType.Message, 512)
	clientPtr := &TcpClient{
		Name:          name,
		OutMsgChannel: msgChannel,
		session:       newSession(conn, msgChannel, make(chan *TcpSession, 1)),
	}

	clientPtr.session.Start()

	return clientPtr, Success
}

func (this *TcpClient) Close() {
	this.session.Stop()
	close(this.OutMsgChannel)
}

// send msg to server
func (this *TcpClient) SendMsg(data []byte) {
	this.session.SendMsg(data)
}
