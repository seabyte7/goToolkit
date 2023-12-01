package tcpLib

import (
	"errors"
	"go.uber.org/zap"
	"goToolkit/logLib"
	"goToolkit/netLib/netType"
	. "goToolkit/protocol"
	"net"
	"sync"
	"time"
)

type TcpServer struct {
	OutMsgChannel chan *netType.ClientServerMsg // Out message channel needs to be referenced by the outside world

	addr           string
	listener       net.Listener
	sessionMap     map[uint32]*TcpSession
	sessionMutex   sync.RWMutex
	disconnectChan chan *TcpSession
	exitChan       chan bool
}

func NewTcpServer(addr string) *TcpServer {
	return &TcpServer{
		OutMsgChannel:  make(chan *netType.ClientServerMsg, 1024),
		addr:           addr,
		listener:       nil,
		sessionMap:     make(map[uint32]*TcpSession, 512),
		disconnectChan: make(chan *TcpSession, 512),
		exitChan:       make(chan bool, 1),
	}
}

func (this *TcpServer) Start() Result {
	listener, err := net.Listen("tcp", this.addr)
	if err != nil {
		logLib.Zap().Error(
			"TcpServer start failed",
			zap.String("addr", this.addr),
			zap.Error(err),
		)
		return err
	}

	this.listener = listener

	go this.acceptThread()
	go this.sessionThread()
	go this.heartBeatThread()

	logLib.Zap().Info("TcpServer start success.", zap.String("addr", this.addr))

	return Success
}

func (this *TcpServer) Stop() {
	close(this.exitChan)

	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}

	for _, item := range this.sessionMap {
		item.Stop()
	}
	this.sessionMap = make(map[uint32]*TcpSession, 8)

	logLib.Zap().Info("TcpServer stop success.", zap.String("addr", this.addr))
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

func (this *TcpServer) removeSessionByList(sessionList []*TcpSession) {
	this.sessionMutex.Lock()
	defer this.sessionMutex.Unlock()

	for _, item := range sessionList {
		delete(this.sessionMap, item.ID)
	}
}

func (this *TcpServer) tryGetSession(sessionID uint32) (*TcpSession, bool) {
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

func (this *TcpServer) heartBeatThread() {
	logLib.Zap().Info("TcpServer heartBeatThread start.")
	willExit := false
	timeTickChan := time.Tick(1 * time.Second)
	for {
		select {
		case <-timeTickChan:
			closeSessionList := make([]*TcpSession, 0, 32)
			this.sessionMutex.Lock()
			for _, item := range this.sessionMap {
				if !item.IsActive() {
					closeSessionList = append(closeSessionList, item)
					continue
				}

				item.HeartBeat()
			}

			// clear not active session
			for _, item := range closeSessionList {
				item.Stop()
				delete(this.sessionMap, item.ID)
			}
			this.sessionMutex.Unlock()
		case <-this.exitChan:
			willExit = true
			break
		}

		if willExit {
			break
		}
	}

	logLib.Zap().Info("TcpServer heartBeatThread exit.")
}
