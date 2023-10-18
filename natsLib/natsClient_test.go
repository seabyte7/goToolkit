package natsLib

import (
	"go.uber.org/zap"
	"goToolkit/logLib"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	logLib.StartLog("./natsLibs_test.log", zap.DebugLevel)
	defer logLib.StopLog()

	natsClientPtr := NewNatsClient("test", "", "", []string{"nats://127.0.0.1:4222"})
	if natsClientPtr == nil {
		t.Error("NewNatsClient failed.")
	}

	natsClientPtr.Close()
}

func TestNatsClient_RegisterMsgHandler(t *testing.T) {
	logLib.StartLog("./natsLibs_test.log", zap.DebugLevel)
	defer logLib.StopLog()

	natsClientPtr := NewNatsClient("test", "", "", []string{"nats://127.0.0.1:4222"})
	if natsClientPtr == nil {
		t.Error("NewNatsClient failed.")
	}
	defer natsClientPtr.Close()

	var api NatsMsgApi = 10001
	natsClientPtr.RegisterMsgHandler(api, func(msg *NatsMsgData) error {
		return nil
	})
}

func TestNatsClient_tryGetMsgHandler(t *testing.T) {
	logLib.StartLog("./natsLibs_test.log", zap.DebugLevel)
	defer logLib.StopLog()

	natsClientPtr := NewNatsClient("test", "", "", []string{"nats://127.0.0.1:4222"})
	if natsClientPtr == nil {
		t.Error("NewNatsClient failed.")
	}
	defer natsClientPtr.Close()

	var api NatsMsgApi = 10001
	natsClientPtr.RegisterMsgHandler(api, func(msg *NatsMsgData) error {
		return nil
	})

	_, exists := natsClientPtr.tryGetMsgHandler(api)
	if !exists {
		t.Error("tryGetMsgHandler failed.")
	}
}

func TestNatsMsgProcessFlow(t *testing.T) {
	logLib.StartLog("./natsLibs_test.log", zap.DebugLevel)
	defer logLib.StopLog()

	natsClientPtr := NewNatsClient("test", "", "", []string{"nats://127.0.0.1:4222"})
	if natsClientPtr == nil {
		t.Error("NewNatsClient failed.")
	}
	defer natsClientPtr.Close()

	err := natsClientPtr.Subscribe("test")
	if err != nil {
		t.Error(err)
	}

	var api NatsMsgApi = 10001
	natsClientPtr.RegisterMsgHandler(api, func(msg *NatsMsgData) error {
		data, _ := msg.Data.([]byte)
		logLib.Zap().Info("TestNatsMsgProcessFlow api receive data", zap.Int("api", int(msg.Api)), zap.String("data", string(data)))

		return nil
	})

	exitChan := make(chan bool, 1)
	go func() {
		select {
		case msg := <-natsClientPtr.MsgChannel:
			natsClientPtr.HandleMsg(msg)
			logLib.Zap().Info("TestNatsMsgProcessFlow HandleMsg msg.")
		case _ = <-exitChan:
			logLib.Zap().Info("TestNatsMsgProcessFlow msgChannel goroutine exit")
			break
		}
	}()

	err = natsClientPtr.Publish("test", api, []byte("hello,nats client."))
	if err != nil {
		t.Error(err)
	}

	logLib.Zap().Info("TestNatsMsgProcessFlow publish msg success.")

	time.Sleep(1 * time.Second)
	close(exitChan)
}
