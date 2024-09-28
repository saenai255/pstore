package pstore

import (
	"bytes"
	"encoding/gob"
)

func serialize(it interface{}) ([]byte, error) {
	writer := new(bytes.Buffer)
	enc := gob.NewEncoder(writer)
	err := enc.Encode(it)

	return writer.Bytes(), err
}

func deserialize(data []byte, out interface{}) error {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	return dec.Decode(out)
}
