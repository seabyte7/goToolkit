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

func BytesToUint8(data []byte) uint8 {
	return data[0]
}

func Uint8ToBytes(num uint8) []byte {
	bytesBuffer := make([]byte, 1)
	bytesBuffer[0] = num
	return bytesBuffer
}
