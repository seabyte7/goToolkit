package encodingLib

import (
	"bytes"
	"encoding/gob"
)

// MarshalGob Serializes the gob
// serialized objects need to invoke the gob. RegisterName/gob. Register to Register
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

// UnmarshalGob deserializes the serialized gob object
func UnmarshalGob(data []byte, to interface{}) error {
	buff := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)
	err := decoder.Decode(to)
	if err != nil {
		// log it
	}

	return err
}
