package natsLib

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap/zapcore"
	"goToolkit/encodingLib"
	"goToolkit/logLib"
)

type NatsClient struct {
	natsConn      *nats.Conn
	msgHandlerMap map[NatsMsgApi]NatsMsgHandler
	MsgChannel    chan *nats.Msg
}

func NewNatsClient(name, user, passwd string, serverList []string) *NatsClient {
	conn, err := connect(name, user, passwd, serverList)
	if err != nil {
		return nil
	}

	return &NatsClient{
		natsConn:      conn,
		msgHandlerMap: make(map[NatsMsgApi]NatsMsgHandler, 64),
		MsgChannel:    make(chan *nats.Msg, 1024*64),
	}
}

func (this *NatsClient) RegisterMsgHandler(api NatsMsgApi, handler NatsMsgHandler) {
	this.msgHandlerMap[api] = handler
}

func (this *NatsClient) tryGetMsgHandler(api NatsMsgApi) (NatsMsgHandler, bool) {
	handler, exists := this.msgHandlerMap[api]
	if !exists {
		return nil, false
	}

	return handler, true
}

func (this *NatsClient) Close() {
	if this.natsConn == nil {
		return
	}

	this.natsConn.Close()
	this.natsConn = nil
}

func (this *NatsClient) HandleMsg(msg *nats.Msg) error {
	var msgData NatsMsgData
	err := encodingLib.UnmarshalGob(msg.Data, &msgData)
	if err != nil {
		logLib.Zap().Error("Unmarshal message failed.",
			zapcore.Field{Key: "method", String: "natsLib.HandleMsg"},
			zapcore.Field{Key: "subject", String: msg.Subject},
			zapcore.Field{Key: "error", String: err.Error()})
		return err
	}

	handler, exists := this.tryGetMsgHandler(msgData.Api)
	if !exists {
		logLib.Zap().Error("Message handler not found.",
			zapcore.Field{Key: "method", String: "natsLib.HandleMsg"},
			zapcore.Field{Key: "subject", String: msg.Subject},
			zapcore.Field{Key: "error", String: "handler not found"})
		return err
	}

	err = handler(msg)

	return err
}

func (this *NatsClient) Publish(subject string, api NatsMsgApi, data interface{}) error {
	publishData := newNatsMsgData(api, data)
	msg, err := encodingLib.MarshalGob(publishData)
	if err != nil {
		return err
	}

	err = this.natsConn.Publish(subject, msg)
	if err != nil {
		logLib.Zap().Error("Publish message failed.",
			zapcore.Field{Key: "method", String: "natsLib.Publish"},
			zapcore.Field{Key: "subject", String: subject},
			zapcore.Field{Key: "error", String: err.Error()})
		return err
	}

	return nil
}

func (this *NatsClient) Subscribe(subject string) error {
	_, err := this.natsConn.ChanSubscribe(subject, this.MsgChannel)
	if err != nil {
		logLib.Zap().Error("Subscribe message failed.",
			zapcore.Field{Key: "method", String: "natsLib.Subscribe"},
			zapcore.Field{Key: "subject", String: subject},
			zapcore.Field{Key: "error", String: err.Error()})
		return err
	}

	return nil
}
