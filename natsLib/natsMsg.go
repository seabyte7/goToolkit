package natsLib

import "github.com/nats-io/nats.go"

type NatsMsgHandler func(*nats.Msg) error

type NatsMsgData struct {
	Api  NatsMsgApi
	Data interface{}
}

func newNatsMsgData(api NatsMsgApi, data interface{}) NatsMsgData {
	return NatsMsgData{
		Api:  api,
		Data: data,
	}
}
