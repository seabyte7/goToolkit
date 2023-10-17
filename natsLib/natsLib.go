package natsLib

var (
	natsClient *NatsClient
)

func Start(name, user, passwd string, serverList []string) {
	natsClient = NewNatsClient(name, user, passwd, serverList)
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
