package natsLib

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

var (
	natsClient *NatsClient
)

func Start(name, user, passwd string, serverList []string) {
	natsClient = NewNatsClient(name, user, passwd, serverList)
	if natsClient == nil {
		panic(fmt.Errorf("NewNatsClient failed"))
	}
}

func Stop() {
	natsClient.Close()
	natsClient = nil
}

func RegisterMsgHandler(api NatsMsgApi, handler NatsMsgHandler) {
	natsClient.RegisterMsgHandler(api, handler)
}

func Publish(subject string, api NatsMsgApi, data []byte) error {
	return natsClient.Publish(subject, api, data)
}

func Subscribe(subject string) error {
	return natsClient.Subscribe(subject)
}

func GetMsgChannel() chan *nats.Msg {
	return natsClient.MsgChannel
}

func HandleMsg(msg *nats.Msg) error {
	return natsClient.HandleMsg(msg)
}
