package tcpLib

import (
	"google.golang.org/protobuf/proto"
	"net"
)

type TcpClient struct {
	Name string

	session *TcpSession
}

func DialTcpServer(name, addr string) (*TcpClient, Result) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	clientPtr := &TcpClient{
		Name:    name,
		session: newClientSession(conn, make(chan *TcpSession, 1)),
	}

	clientPtr.session.Start()

	return clientPtr, Success
}

func (this *TcpClient) Close() {
	this.session.Stop()
}

// send msg to server
func (this *TcpClient) SendMsg(data []byte) {
	this.session.SendMsg(data)
}

func (this *TcpClient) SendPBMsg(msg proto.Message) Result {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	this.session.SendMsg(data)

	return Success
}

func (this *TcpClient) GetReceiveMsgChan() chan *ClientServerMsg {
	return this.session.ReceiveMsgChan
}
