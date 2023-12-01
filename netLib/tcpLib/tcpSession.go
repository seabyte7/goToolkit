package tcpLib

import (
	"go.uber.org/zap"
	"goToolkit/logLib"
	"goToolkit/netLib/netType"
	"io"
	"net"
	"sync"
	"time"
)

const (
	MaxHeartBeatIntervalMs = 5 * 1000 // 60s
)

type TcpSession struct {
	ID uint32

	SendMsgChannel     chan *netType.ClientServerMsg
	OutMsgChanRef      chan *netType.ClientServerMsg // Referenced external out message channel
	disconnectChanRef  chan *TcpSession              // Referenced external disconnected channel
	readWriteCloseChan chan struct{}                 // The internal read/write channel is closed

	conn        net.Conn
	recvDataBuf []byte // recv data buffer

	recvBytes uint64 // recv bytes
	sendBytes uint64 // send bytes

	lastActiveTimeMs int64 // last active time

	exitChan  chan struct{}
	closeOnce sync.Once
}

func newSession(conn net.Conn, outMsgChan chan *netType.ClientServerMsg, disconnectChan chan *TcpSession) *TcpSession {
	return &TcpSession{
		ID:                 acquireID(),
		SendMsgChannel:     make(chan *netType.ClientServerMsg, 512),
		OutMsgChanRef:      outMsgChan,
		disconnectChanRef:  disconnectChan,
		readWriteCloseChan: make(chan struct{}, 2),
		conn:               conn,
		recvDataBuf:        make([]byte, 0),
		recvBytes:          0,
		sendBytes:          0,
		lastActiveTimeMs:   time.Now().UnixMilli(),
		exitChan:           make(chan struct{}, 1),
	}
}

func (this *TcpSession) Start() {
	logLib.Zap().Info("session start", zap.Uint32("sessionID", this.ID))
	go this.monitorThread()
	go this.recvThread()
	go this.sendThread()
}

func (this *TcpSession) Stop() {
	this.closeOnce.Do(func() {
		logLib.Zap().Info("session stop", zap.Uint32("sessionID", this.ID))
		close(this.exitChan)
		if this.conn != nil {
			this.conn.Close()
			this.conn = nil
		}
	})
}

func (this *TcpSession) monitorThread() {
	select {
	case <-this.exitChan:
		return
	case <-this.readWriteCloseChan:
		this.disconnectChanRef <- this
		return
	}
}

func (this *TcpSession) recvThread() {
	logLib.Zap().Info("session recv thread start", zap.Uint32("sessionID", this.ID))
	defer func() {
		logLib.Zap().Info("session recv thread stop", zap.Uint32("sessionID", this.ID))
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
				logLib.Zap().Error("session:%d read 0 byte", zap.Uint32("sessionID", this.ID))
				return
			}

			this.recvDataBuf = append(this.recvDataBuf, recvBuf[:count]...)
			this.recvBytes += uint64(count)
			this.lastActiveTimeMs = time.Now().UnixMilli()

			msgList := this.parseMessage()
			for _, msg := range msgList {
				switch msg.Cmd {
				case netType.CmdHeartbeat:
					this.OnRecvHeartBeat()
				case netType.CmdMsg:
					this.OutMsgChanRef <- msg
				default:
					logLib.Sugar().Errorf("session:%d recv unknown cmd:%d", this.ID, msg.Cmd) // todo consider use zap
				}
			}
		}
	}

	logLib.Zap().Info("session recv thread exit", zap.Uint32("sessionID", this.ID))
}

// parse message packet
func (this *TcpSession) parseMessage() []*netType.ClientServerMsg {
	msgList := make([]*netType.ClientServerMsg, 0, 4)
	maxTryCount := 4

	for i := 0; i < maxTryCount; i++ {
		if len(this.recvDataBuf) == 0 {
			break
		}

		msgLen := netType.BytesToUint32(this.recvDataBuf[:netType.MsgLengthLen])
		bufLen := len(this.recvDataBuf)
		if msgLen < uint32(bufLen) {
			return msgList
		}
		cmd := netType.BytesToUint8(this.recvDataBuf[netType.MsgLengthLen : netType.MsgLengthLen+netType.MsgCmdLen])

		msg := netType.NewClientServerMsg(cmd, this.recvDataBuf[netType.MsgHeaderLen:msgLen])
		msgList = append(msgList, msg)

		this.recvDataBuf = this.recvDataBuf[msgLen:]
	}

	return msgList
}

func (this *TcpSession) sendThread() {
	logLib.Zap().Info("session send thread start", zap.Uint32("sessionID", this.ID))
	defer func() {
		logLib.Zap().Info("session send thread stop", zap.Uint32("sessionID", this.ID))
		this.readWriteCloseChan <- struct{}{}
	}()

	for {
		willExit := false
		select {
		case <-this.exitChan:
			willExit = true
			break
		case msg := <-this.SendMsgChannel:
			count, err := this.conn.Write(msg.ToBytes())
			if err != nil {
				logLib.Sugar().Errorf("session:%d write error:%v", this.ID, err)
				break
			}
			logLib.Sugar().Debugf("session:%d write msg:%v, length:%d", this.ID, msg, count)
			this.sendBytes += uint64(count)
			this.lastActiveTimeMs = time.Now().UnixMilli()
		}

		if willExit {
			break
		}
	}

	logLib.Zap().Info("session send thread exit", zap.Uint32("sessionID", this.ID))
}

func (this *TcpSession) SendHeartBeatMsg() {
	this.sendMsgByCmd(netType.CmdHeartbeat, []byte{})
}

func (this *TcpSession) SendMsg(data []byte) {
	this.sendMsgByCmd(netType.CmdMsg, data)
}

func (this *TcpSession) sendMsgByCmd(cmd uint8, data []byte) {
	msg := netType.NewClientServerMsg(cmd, data)

	this.SendMsgChannel <- msg
}

func (this *TcpSession) HeartBeat() {
	this.SendHeartBeatMsg()
}

func (this *TcpSession) OnRecvHeartBeat() {
	this.lastActiveTimeMs = time.Now().UnixMilli()
	logLib.Sugar().Debugf("session:%d recv heartbeat", this.ID)
}

func (this *TcpSession) IsActive() bool {
	nowTime := time.Now()
	if nowTime.UnixMilli()-this.lastActiveTimeMs > MaxHeartBeatIntervalMs {
		return false
	}

	return true
}
