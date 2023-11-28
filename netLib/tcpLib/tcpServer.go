package tcpLib

import (
	"goToolkit/logLib"
	"net"
	"sync"
)

type TcpServer struct {
	Addr string

	conn net.Conn

	CloseOnce sync.Once
	exitChan  chan bool
}

func NewTcpServer(addr string) *TcpServer {
	return &TcpServer{
		Addr: addr,
		//CloseOnce:
		exitChan: make(chan bool, 1),
	}
}

func (this *TcpServer) Start() {
	net.Listen("tcp", this.Addr)
}

func (this *TcpServer) Stop() {
	this.CloseOnce.Do(func() {
		close(this.exitChan)
	})
}

func (this *TcpServer) acceptThread() {
	logLib.Zap().Info("TcpServer acceptThread start.")

	willExit := false
	for {

		select {
		case <-this.exitChan:
			willExit = true

		}
	}
}
