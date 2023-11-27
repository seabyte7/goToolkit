package tcpLib

import (
	"goToolkit/logLib"
	"goToolkit/netLib/netType"
	. "goToolkit/protocol"
	"io"
	"net"
	"sync"
)

type TcpClient struct {
	SessionID uint64

	RecvMsgChannel chan *netType.Message
	SendMsgChannel chan *netType.Message

	conn        net.Conn
	recvDataBuf []byte // recv data buffer

	exitChan  chan struct{}
	closeOnce sync.Once
}

func DialTcpServer(addr string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	clientPtr := &TcpClient{
		SessionID:      acquireID(),
		RecvMsgChannel: make(chan *netType.Message, 512),
		SendMsgChannel: make(chan *netType.Message, 512),
		conn:           conn,
		recvDataBuf:    make([]byte, 0),
		exitChan:       make(chan struct{}, 1),
	}

	go clientPtr.recvThread()
	go clientPtr.sendThread()

	return clientPtr, Success
}

func (this *TcpClient) Close() {
	this.closeOnce.Do(func() {
		close(this.exitChan)
		this.conn.Close()
	})
}

func (this *TcpClient) recvThread() {
	for {
		recvBuf := make([]byte, 0, 1024)
		count, err := this.conn.Read(recvBuf)
		if err != nil {
			if err == io.EOF {
				logLib.Sugar().Infof("session:%d read eof, socket closed.", this.SessionID)
				break
			}
			logLib.Sugar().Infof("session:%d read error:%v", this.SessionID, err)
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
func (this *TcpClient) parseMessage() []*netType.Message {
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

func (this *TcpClient) sendThread() {
	for {
		msg := <-this.SendMsgChannel
		count, err := this.conn.Write(msg.ToMsgBytes())
		if err != nil {
			logLib.Sugar().Errorf("session:%d write error:%v", this.SessionID, err)
			break
		}
		logLib.Sugar().Debugf("session:%d write msg:%v, length:%d", this.SessionID, msg, count)
	}
}

// send msg to server
func (this *TcpClient) SendMsg(data []byte) {
	msgLen := netType.MsgHeaderLen + uint32(len(data))
	msg := netType.NewMessage(msgLen, data)

	this.SendMsgChannel <- msg
}
