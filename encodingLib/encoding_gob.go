package encodingLib

import (
	"bytes"
	"encoding/gob"
)

// MarshalGob 序列化gob
// 序列化的对象需要先调用gob.RegisterName/gob.Register注册
func MarshalGob(data interface{}) ([]byte, error) {
	buff := &bytes.Buffer{}
	encoder := gob.NewEncoder(buff)

	err := encoder.Encode(data)
	if err != nil {
		// log it
		return nil, err
	}

	return buff.Bytes(), nil
}

// UnmarshalGob 反序列化gob
func UnmarshalGob(data []byte, to interface{}) {
	buff := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)
	err := decoder.Decode(to)
	if err != nil {
		// log it
		return
	}
}
