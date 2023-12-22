package main

import (
	"fmt"
	"github.com/seabyte7/goToolkitnetLib/tcpLib"
	"math/rand"
	"time"
)

var (
	tcpServer         *tcpLib.TcpServer
	stopTcpServerChan = make(chan bool, 1)
)

func tcpLibTestFunc() {
	myRandPtr := rand.New(rand.NewSource(time.Now().Unix()))
	addr := fmt.Sprintf(":%d", 30000+myRandPtr.Int31n(10000))
	StartTcpServer(addr)
	defer StopTcpServer()
	go processTcpServerMsgThread()

	clientPtr, result := DialogTcpServer(addr)
	if result != Success {
		panic(fmt.Errorf("DialogTcpServer failed, err:%v", result))
	}
	defer clientPtr.Close()
	go processTcpClientMsgThread(clientPtr)

	for i := 0; i < 5; i++ {
		clientPtr.SendMsg([]byte(fmt.Sprintf("hello world %d", i)))
		time.Sleep(1 * time.Second)
	}
}

func StartTcpServer(addr string) {
	tcpServer = tcpLib.NewTcpServer(addr)

	result := tcpServer.Start()
	if result != Success {
		panic(fmt.Errorf("StartTcpServer NewTcpServer failed, err:%v", result))
	}
}

func StopTcpServer() {
	close(stopTcpServerChan)
	tcpServer.Stop()
	tcpServer = nil
}

// go test -v -run TestDialogTcpServer
// client
func DialogTcpServer(addr string) (*tcpLib.TcpClient, Result) {
	tcpClientPtr, err := tcpLib.DialTcpServer("test", addr)
	if err != nil {
		return nil, err
	}

	return tcpClientPtr, Success
}

func processTcpServerMsgThread() {
	for {
		select {
		case <-stopTcpServerChan:
			return
		case msg := <-tcpServer.ReceiveMsgChan:
			fmt.Printf("TcpServer recv msg:%s\n", string(msg.Data))
			msg.Session.SendMsg([]byte(fmt.Sprintf("hello, client:%v echo:%v", msg.Session.ID, string(msg.Data))))
		case sessionPtr := <-tcpServer.AcceptSessionChan:
			sessionPtr.SendMsg([]byte(fmt.Sprintf("hello, client:%v", sessionPtr.ID)))
		}
	}
}

func processTcpClientMsgThread(clientPtr *tcpLib.TcpClient) {
	for {
		select {
		case msg := <-clientPtr.GetReceiveMsgChan():
			fmt.Printf("client receive msg:%v\n", string(msg.Data))
		}
	}
}
