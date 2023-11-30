package tcpLib

import (
	"errors"
	"go.uber.org/zap"
	"goToolkit/logLib"
	"goToolkit/netLib/netType"
	. "goToolkit/protocol"
	"net"
	"sync"
)

type TcpServer struct {
	Addr string

	listener     net.Listener
	sessionMap   map[uint64]*TcpSession
	sessionMutex sync.RWMutex

	OutMsgChannel  chan *netType.Message
	disconnectChan chan *TcpSession

	CloseOnce sync.Once
	exitChan  chan bool
}

func NewTcpServer(addr string) *TcpServer {
	return &TcpServer{
		Addr:           addr,
		listener:       nil,
		sessionMap:     make(map[uint64]*TcpSession, 512),
		OutMsgChannel:  make(chan *netType.Message, 1024),
		disconnectChan: make(chan *TcpSession, 512),
		//CloseOnce:
		exitChan: make(chan bool, 1),
	}
}

func (this *TcpServer) Start() Result {
	listener, err := net.Listen("tcp", this.Addr)
	if err != nil {
		logLib.Zap().Error(
			"TcpServer start failed",
			zap.String("addr", this.Addr),
			zap.Error(err),
		)
		return err
	}

	this.listener = listener

	go this.acceptThread()
	go this.sessionThread()

	logLib.Zap().Info("TcpServer start success.", zap.String("addr", this.Addr))

	return Success
}

func (this *TcpServer) Stop() {
	this.listener.Close()
	this.listener = nil

	for _, item := range this.sessionMap {
		item.Stop()
	}

	logLib.Zap().Info("TcpServer stop success.", zap.String("addr", this.Addr))
}

func (this *TcpServer) addSession(session *TcpSession) {
	this.sessionMutex.Lock()
	defer this.sessionMutex.Unlock()

	this.sessionMap[session.ID] = session
}

func (this *TcpServer) removeSession(session *TcpSession) {
	this.sessionMutex.Lock()
	defer this.sessionMutex.Unlock()

	delete(this.sessionMap, session.ID)
}

func (this *TcpServer) tryGetSession(sessionID uint64) (*TcpSession, bool) {
	this.sessionMutex.RLock()
	defer this.sessionMutex.RUnlock()

	sessionPtr, ok := this.sessionMap[sessionID]
	if !ok {
		return nil, false
	}

	return sessionPtr, true
}

func (this *TcpServer) acceptThread() {
	logLib.Zap().Info("TcpServer acceptThread start.")
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logLib.Zap().Info("TcpServer acceptThread exit.")
				return
			}

			logLib.Zap().Error("TcpServer accept failed", zap.Error(err))
			return
		}

		sessionPtr := newSession(conn, this.OutMsgChannel, this.disconnectChan)
		sessionPtr.Start()
		this.addSession(sessionPtr)
		logLib.Sugar().Infof("New session:%d connected.", sessionPtr.ID)
	}

	logLib.Zap().Info("TcpServer acceptThread exit.")
}

func (this *TcpServer) sessionThread() {
	logLib.Zap().Info("TcpServer sessionThread start.")
	willExit := false
	for {
		select {
		case <-this.exitChan:
			willExit = true
			break
		case sessionPtr := <-this.disconnectChan:
			sessionPtr.Stop()
			this.removeSession(sessionPtr)
		}

		if willExit {
			break
		}
	}

	logLib.Zap().Info("TcpServer sessionThread exit.")
}
