package tcpLib

import "sync/atomic"

var (
	autoID uint32
)

func acquireID() uint32 {
	return atomic.AddUint32(&autoID, 1)
}
