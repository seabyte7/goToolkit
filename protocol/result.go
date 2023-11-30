package protocol

import "errors"

type Result error

var (
	Success Result = nil
)

// 错误码
var (
	RecvDataNotEnough Result = errors.New("RecvDataNotEnough")
	NetListenFailed   Result = errors.New("NetListenFailed")
)
