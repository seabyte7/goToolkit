package netType

import "encoding/binary"

var (
	NetByteOrder = binary.BigEndian
)

func Uint32ToBytes(num uint32) []byte {
	bytesBuffer := make([]byte, 4)
	NetByteOrder.PutUint32(bytesBuffer, num)
	return bytesBuffer
}

func BytesToUint32(data []byte) uint32 {
	return NetByteOrder.Uint32(data)
}
