package tcpLib

import (
	"go.uber.org/zap"
	"goToolkit/logLib"
	"goToolkit/netLib/netType"
	"io"
	"net"
	"sync"
)

type TcpSession struct {
	ID uint64

	SendMsgChannel     chan *netType.Message
	OutMsgChanRef      chan *netType.Message // Referenced external out message channel
	disconnectChanRef  chan *TcpSession      // Referenced external disconnected channel
	readWriteCloseChan chan struct{}         // The internal read/write channel is closed

	conn        net.Conn
	recvDataBuf []byte // recv data buffer

	exitChan  chan struct{}
	closeOnce sync.Once
}

func newSession(conn net.Conn, outMsgChan chan *netType.Message, disconnectChan chan *TcpSession) *TcpSession {
	return &TcpSession{
		ID:                 acquireID(),
		SendMsgChannel:     make(chan *netType.Message, 512),
		OutMsgChanRef:      outMsgChan,
		disconnectChanRef:  disconnectChan,
		readWriteCloseChan: make(chan struct{}, 2),
		conn:               conn,
		recvDataBuf:        make([]byte, 0),
		exitChan:           make(chan struct{}, 1),
	}
}

func (this *TcpSession) Start() {
	logLib.Zap().Info("session start", zap.Uint64("sessionID", this.ID))
	go this.logicThread()
	go this.recvThread()
	go this.sendThread()
}

func (this *TcpSession) Stop() {
	this.closeOnce.Do(func() {
		logLib.Zap().Info("session stop", zap.Uint64("sessionID", this.ID))
		close(this.exitChan)
		if this.conn != nil {
			this.conn.Close()
			this.conn = nil
		}
	})
}

func (this *TcpSession) logicThread() {
	select {
	case <-this.exitChan:
		return
	case <-this.readWriteCloseChan:
		this.disconnectChanRef <- this
		return
	}
}

func (this *TcpSession) recvThread() {
	logLib.Zap().Info("session recv thread start", zap.Uint64("sessionID", this.ID))
	defer func() {
		logLib.Zap().Info("session recv thread stop", zap.Uint64("sessionID", this.ID))
		this.readWriteCloseChan <- struct{}{}
	}()

	for {
		select {
		case <-this.exitChan:
			break
		default:
			recvBuf := make([]byte, 1024)
			count, err := this.conn.Read(recvBuf)
			if err != nil {
				if err == io.EOF {
					logLib.Sugar().Infof("session:%d read eof, socket closed.", this.ID)
					return
				}
				logLib.Sugar().Infof("session:%d read error:%v", this.ID, err)
				return
			}
			if count == 0 {
				// If the caller wanted a zero byte read, return immediately
				// without trying (but after acquiring the readLock).
				// Otherwise syscall.Read returns 0, nil which looks like
				// io.EOF.
				// TODO(bradfitz): make it wait for readability? (Issue 15735)
				// src/internal/poll/fd_unix.go
				logLib.Zap().Error("session:%d read 0 byte", zap.Uint64("sessionID", this.ID))
				return
			}

			this.recvDataBuf = append(this.recvDataBuf, recvBuf[:count]...)

			msgList := this.parseMessage()
			for _, msg := range msgList {
				this.OutMsgChanRef <- msg
			}
		}
	}
}

// parse message packet
func (this *TcpSession) parseMessage() []*netType.Message {
	msgList := make([]*netType.Message, 0, 4)
	maxTryCount := 4

	for i := 0; i < maxTryCount; i++ {
		if len(this.recvDataBuf) == 0 {
			break
		}

		msgLen := netType.BytesToUint32(this.recvDataBuf[:netType.MsgHeaderLen])
		bufLen := len(this.recvDataBuf)
		if msgLen < uint32(bufLen) {
			return msgList
		}

		msg := netType.NewMessage(msgLen, this.recvDataBuf[netType.MsgHeaderLen:msgLen])
		msgList = append(msgList, msg)

		this.recvDataBuf = this.recvDataBuf[msgLen:]
	}

	return msgList
}

func (this *TcpSession) sendThread() {
	logLib.Zap().Info("session send thread start", zap.Uint64("sessionID", this.ID))
	defer func() {
		logLib.Zap().Info("session send thread stop", zap.Uint64("sessionID", this.ID))
		this.readWriteCloseChan <- struct{}{}
	}()

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

func (this *TcpSession) SendMsg(data []byte) {
	msgLen := netType.MsgHeaderLen + uint32(len(data))
	msg := netType.NewMessage(msgLen, data)

	this.SendMsgChannel <- msg
}
