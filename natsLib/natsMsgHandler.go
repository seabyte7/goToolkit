package natsLib

import (
	"github.com/nats-io/nats.go"
	"log"
)

type NatsMsgHandler func(*nats.Msg) error

func handleMsg(msg *nats.Msg) {
	log.Println("Received a message: ", string(msg.Data))
}

var (
	msgHandlerMap = make(map[string]NatsMsgHandler, 64)
)

func RegisterMsgHandler(subject string, handler NatsMsgHandler) {
	msgHandlerMap[subject] = handler
}

func tryGetMsgHandler(subject string) (NatsMsgHandler, bool) {
	handler, exists := msgHandlerMap[subject]
	if !exists {
		return nil, false
	}

	return handler, true
}
