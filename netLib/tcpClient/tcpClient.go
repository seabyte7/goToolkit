package tcpClient

import (
	"goToolkit/netLib/netType"
	"net"
	"sync/atomic"
)

var (
	autoID int64
)

type TcpClient struct {
	conn          net.Conn
	SessionId     int64
	recvDataBuf   []byte
	InMsgChannel  chan netType.Message
	OutMsgChannel chan netType.Message

	exitChan chan struct{}
}

func AcquireID() int64 {
	return atomic.AddInt64(&autoID, 1)
}

func DialTcpServer(addr string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client := &TcpClient{
		conn:          conn,
		SessionId:     AcquireID(),
		recvDataBuf:   make([]byte, 0),
		InMsgChannel:  make(chan netType.Message, 100),
		OutMsgChannel: make(chan netType.Message, 100),
	}

	go client.recvThread()
	go client.send()

	return client, nil
}

func (this *TcpClient) Close() {
	close(this.exitChan)
	this.conn.Close()
}

func (this *TcpClient) recvThread() {
	for {

		this.conn.Read()
		msg, err := netType.ReadMessage(this.conn, this.recvDataBuf)
		if err != nil {
			return
		}

		this.InMsgChannel <- msg
	}
}

func (this *TcpClient) sendThread() {
	for {
		msg := <-this.OutMsgChannel
		netType.WriteMessage(this.conn, msg)
	}
}
