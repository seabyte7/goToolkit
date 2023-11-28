package tcpLib

import (
	"goToolkit/logLib"
	"goToolkit/netLib/netType"
	"io"
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

func (this *Session) recvThread() {
	for {
		recvBuf := make([]byte, 0, 1024)
		count, err := this.conn.Read(recvBuf)
		if err != nil {
			if err == io.EOF {
				logLib.Sugar().Infof("session:%d read eof, socket closed.", this.ID)
				break
			}
			logLib.Sugar().Infof("session:%d read error:%v", this.ID, err)
			break
		}

		this.recvDataBuf = append(this.recvDataBuf, recvBuf[:count]...)

		msgList := this.parseMessage()
		for _, msg := range msgList {
			this.RecvMsgChannel <- msg
		}
	}
}

// parse message packet
func (this *Session) parseMessage() []*netType.Message {
	msgList := make([]*netType.Message, 0, 4)
	maxTryCount := 4

	for i := 0; i < maxTryCount; i++ {
		msgLen := netType.BytesToUint32(this.recvDataBuf[:netType.MsgHeaderLen])
		bufLen := len(this.recvDataBuf)
		if msgLen < uint32(bufLen) {
			return msgList
		}

		msg := netType.NewMessage(msgLen, this.recvDataBuf[:msgLen])
		msgList = append(msgList, msg)

		this.recvDataBuf = this.recvDataBuf[msgLen:]
	}

	return msgList
}

func (this *Session) sendThread() {
	for {
		msg := <-this.SendMsgChannel
		count, err := this.conn.Write(msg.ToMsgBytes())
		if err != nil {
			logLib.Sugar().Errorf("session:%d write error:%v", this.ID, err)
			break
		}
		logLib.Sugar().Debugf("session:%d write msg:%v, length:%d", this.ID, msg, count)
	}
}

// send msg to server
func (this *Session) SendMsg(data []byte) {
	msgLen := netType.MsgHeaderLen + uint32(len(data))
	msg := netType.NewMessage(msgLen, data)

	this.SendMsgChannel <- msg
}
