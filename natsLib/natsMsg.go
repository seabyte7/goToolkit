package natsLib

type NatsMsgHandler func(*NatsMsgData) error

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
