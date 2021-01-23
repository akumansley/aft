package oplog

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type Encoder func(interface{}) ([]byte, error)
type Decoder func([]byte) (interface{}, error)

type codecLog struct {
	LogStore
	encode Encoder
	decode Decoder
}

func NewLog(store LogStore, enc Encoder, dec Decoder) OpLog {
	return &codecLog{LogStore: store, encode: enc, decode: dec}
}

func (c *codecLog) Iterator() Iterator {
	bi := c.LogStore.Iterator()
	return &codecIter{ByteIterator: bi, log: c}
}

func (c *codecLog) Scan(count, offset int) (result []interface{}, err error) {
	result = []interface{}{}

	vals, err := c.LogStore.Scan(count, offset)
	if err != nil {
		return
	}
	for _, v := range vals {
		var decoded interface{}
		decoded, err = c.decode(v)
		if err != nil {
			return
		}
		result = append(result, decoded)
	}
	return
}

func (c *codecLog) Log(object interface{}) error {
	bytes, err := c.encode(object)
	if err != nil {
		return err
	}
	return c.LogStore.Log(bytes)
}

type codecIter struct {
	ByteIterator
	log *codecLog
	val interface{}
	err error
}

func GobDecoder(encoded []byte) (interface{}, error) {
	buf := bytes.NewBuffer(encoded)
	var val interface{}
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&val)
	return val, err
}

func GobEncoder(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&v)
	return buf.Bytes(), err
}

func BytesDecoder(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func BytesEncoder(v interface{}) (bts []byte, err error) {
	bts, ok := v.([]byte)
	if !ok {
		err = errors.New("BytesEncoder expected []byte")
	}
	return bts, err
}

func GobLog(logStore LogStore) OpLog {
	return &codecLog{LogStore: logStore, encode: GobEncoder, decode: GobDecoder}
}

func BinaryLog(logStore LogStore) OpLog {
	return &codecLog{LogStore: logStore, encode: BytesEncoder, decode: BytesDecoder}
}

func (c *codecIter) Next() bool {
	ok := c.ByteIterator.Next()
	if !ok {
		c.err = c.ByteIterator.Err()
		return false
	} else {
		val, err := c.log.decode(c.ByteIterator.Value())
		if err != nil {
			c.err = err
			return false
		}
		c.val = val
		return true
	}
}

func (c *codecIter) Value() interface{} {
	return c.val
}

func (c *codecIter) Err() error {
	return c.err
}
