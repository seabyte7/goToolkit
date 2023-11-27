package tcpLib

import "sync/atomic"

var (
	autoID uint64
)

func acquireID() uint64 {
	return atomic.AddUint64(&autoID, 1)
}
