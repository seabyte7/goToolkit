package tcpLib

//func TestNewTcpServer(t *testing.T) {
//	t.Log("TestNewTcpServer")
//	myRandPtr := rand.New(rand.NewSource(time.Now().Unix()))
//
//	addr := fmt.Sprintf(":%d", 30000+myRandPtr.Int31n(10000))
//	tcpServer := NewTcpServer(addr)
//
//	result := tcpServer.Start()
//	if result != Success {
//		t.Errorf("TestNewTcpServer NewTcpServer failed, err:%v", result)
//		return
//	}
//
//	defer tcpServer.Stop()
//
//	t.Log("TestNewTcpServer success")
//}

//func TestDialogTcpServer(t *testing.T) {
//	t.Log("TestDialogTcpServer")
//	myRandPtr := rand.New(rand.NewSource(time.Now().Unix()))
//
//	addr := fmt.Sprintf(":%d", 30000+myRandPtr.Int31n(10000))
//	tcpServer := NewTcpServer(addr)
//
//	result := tcpServer.Start()
//	if result != Success {
//		t.Errorf("TestDialogTcpServer NewTcpServer failed, err:%v", result)
//		return
//	}
//
//	defer tcpServer.Stop()
//
//	time.Sleep(5 * time.Second)
//
//	tcpClient, err := DialTcpServer("test", addr)
//	if err != nil {
//		t.Errorf("TestDialogTcpServer DialTcpServer failed, err:%v", err)
//		return
//	}
//
//	defer tcpClient.Close()
//
//	time.Sleep(1 * time.Second)
//
//	t.Log("TestDialogTcpServer success")
//}
