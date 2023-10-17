package natsLib

type NatsMsgApi int32

var (
	NatsMsgApiUnknown NatsMsgApi = 0
)

type NatsMsgData struct {
	Api  NatsMsgApi
	Data interface{}
}

func NewNatsMsgData(api NatsMsgApi, data interface{}) NatsMsgData {
	return NatsMsgData{
		Api:  api,
		Data: data,
	}
}
