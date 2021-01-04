package db

import (
	"bytes"

	iradix "github.com/hashicorp/go-immutable-radix"
)

type snapshot interface {
	PrefixIterator(prefix []byte) KVIterator
	Get(key []byte) (interface{}, error)
}

type kv interface {
	PrefixIterator(prefix []byte) KVIterator
	Get(key []byte) (interface{}, error)
	Put(key []byte, value interface{}) error
	Delete(key []byte) error
	Snapshot() snapshot
}

// iradixKV

func newIRadixKV() kv {
	return &iradixKV{t: iradix.New()}
}

type iradixKV struct {
	t *iradix.Tree
}

func (kv *iradixKV) PrefixIterator(prefix []byte) KVIterator {
	it := kv.t.Root().Iterator()
	it.SeekPrefix(prefix)
	return &prefixiter{it, nil, nil, nil, prefix}
}

func (kv *iradixKV) Get(key []byte) (interface{}, error) {
	val, ok := kv.t.Get(key)
	if !ok {
		return nil, ErrNotFound
	}
	return val, nil
}

func (kv *iradixKV) Put(key []byte, value interface{}) error {
	newT, _, _ := kv.t.Insert(key, value)
	kv.t = newT
	return nil
}

func (kv *iradixKV) Delete(key []byte) error {
	newT, _, ok := kv.t.Delete(key)
	if !ok {
		return ErrNotFound
	}
	kv.t = newT
	return nil
}

func (kv *iradixKV) Snapshot() snapshot {
	return &iradixKV{kv.t}
}

// prefixiter

type prefixiter struct {
	it     *iradix.Iterator
	val    interface{}
	key    []byte
	err    error
	prefix []byte
}

func (i *prefixiter) Value() interface{} {
	if i.err != nil {
		panic("Called value after err")
	}
	return i.val
}

func (i *prefixiter) Key() []byte {
	if i.err != nil {
		panic("Called value after err")
	}
	return i.key
}

func (i *prefixiter) Entry() ([]byte, interface{}) {
	if i.err != nil {
		panic("Called value after err")
	}
	return i.key, i.val
}

func (i *prefixiter) Err() error {
	return i.err
}

func (i *prefixiter) Next() bool {
	for k, v, ok := i.it.Next(); ok; k, v, ok = i.it.Next() {
		if !bytes.HasPrefix(k, i.prefix) {
			i.err = Done
			return false
		}

		i.val = v
		i.key = k
		return true
	}
	i.err = Done
	return false
}
