package natsLib

import (
	"github.com/seabyte7/goToolkit/logLib"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestNatsStartAndStop(t *testing.T) {
	logLib.StartLog("./natsLibs_test.log", zapcore.DebugLevel)
	defer logLib.StopLog()

	Start("test", "", "", []string{"nats://127.0.0.1:4222"})
	defer Stop()
}

func TestGlobalMsgProcessFlowHandler(t *testing.T) {
	logLib.StartLog("./natsLibs_test.log", zapcore.DebugLevel)
	defer logLib.StopLog()

	Start("test", "", "", []string{"nats://127.0.0.1:4222"})
	defer Stop()

	err := Subscribe("test")
	if err != nil {
		t.Error(err)
	}

	var api NatsMsgApi = 10001
	RegisterMsgHandler(api, func(msg *NatsMsgData) error {
		data, _ := msg.Data.([]byte)
		logLib.Sugar().Debug("TestGlobalMsgProcessFlowHandler api receive data.", msg.Api, string(data))
		return nil
	})

	exitChan := make(chan bool, 1)
	go func() {
		select {
		case msg := <-GetMsgChannel():
			HandleMsg(msg)
			logLib.Sugar().Debug("TestGlobalMsgProcessFlowHandler HandleMsg msg.")
		case _ = <-exitChan:
			logLib.Sugar().Debug("TestGlobalMsgProcessFlowHandler msgChannel goroutine exit")
			break
		}
	}()

	err = Publish("test", api, []byte("hello,nats client."))
	if err != nil {
		t.Error(err)
	}
	logLib.Sugar().Debugf("TestGlobalMsgProcessFlowHandler publish msg success.")

	time.Sleep(1 * time.Second)
	close(exitChan)
}
