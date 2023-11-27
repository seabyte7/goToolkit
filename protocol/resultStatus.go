package protocol

import "errors"

type ResultStatus error

var (
	Success ResultStatus = nil
)

// 错误码
var (
	RecvDataNotEnough ResultStatus = errors.New("RecvDataNotEnough")
)
