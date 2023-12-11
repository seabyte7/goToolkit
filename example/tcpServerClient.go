package main

import (
	"fmt"
	"goToolkit/netLib/tcpLib"
	"math/rand"
	"time"

	. "goToolkit/protocol"
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
	go processTcpServerMsg()

	clientPtr, result := DialogTcpServer(addr)
	if result != Success {
		panic(fmt.Errorf("DialogTcpServer failed, err:%v", result))
	}
	defer clientPtr.Close()

	for i := 0; i < 5; i++ {
		clientPtr.SendMsg([]byte(fmt.Sprintf("hello world %d", i)))
		time.Sleep(3 * time.Second)
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

func processTcpServerMsg() {
	for {
		select {
		case <-stopTcpServerChan:
			return
		case msg := <-tcpServer.OutMsgChannel:
			fmt.Printf("TcpServer recv msg:%s\n", string(msg.Data))
		}
	}
}
