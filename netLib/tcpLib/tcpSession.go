package tcpLib

import (
	"goToolkit/netLib/netType"
	"net"
	"sync"
)

type Session struct {
	ID uint64

	RecvMsgChannel chan *netType.Message
	SendMsgChannel chan *netType.Message

	conn        net.Conn
	recvDataBuf []byte // recv data buffer

	exitChan  chan struct{}
	closeOnce sync.Once
}

func newSession(conn net.Conn) *Session {
	return &Session{
		ID:             acquireID(),
		RecvMsgChannel: make(chan *netType.Message, 512),
		SendMsgChannel: make(chan *netType.Message, 512),
		conn:           conn,
		recvDataBuf:    make([]byte, 0),
		exitChan:       make(chan struct{}, 1),
	}
}
