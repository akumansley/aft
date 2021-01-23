package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"awans.org/aft/internal/oplog"
	"github.com/google/uuid"
)

func DBOpLog(builder *Builder, store oplog.LogStore) oplog.OpLog {
	encoder := func(entry interface{}) ([]byte, error) {
		ops := entry.(Ops)
		data, err := ops.Serialize()
		return data, err
	}
	decoder := func(entry []byte) (interface{}, error) {
		var ops Ops
		err := ops.Deserialize(builder, entry)
		return ops, err
	}
	return oplog.NewLog(store, encoder, decoder)
}

func DBFromLog(appDB DB, l oplog.OpLog) error {
	iter := l.Iterator()
	for iter.Next() {
		rwtx := appDB.NewRWTx()
		var ops Ops
		v := iter.Value()
		ops, ok := v.(Ops)
		if !ok {
			return errors.New("Expected Ops from db log")
		}
		for _, op := range []Operation(ops) {
			op.Replay(rwtx)
		}
		err := rwtx.Commit()
		if err != nil {
			return err
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	return nil
}

func MakeTransactionLogger(l oplog.OpLog) func(BeforeCommit) {
	logger := func(event BeforeCommit) {
		opslice := event.Tx.Operations()

		if len(opslice) == 0 {
			return
		}
		ops := Ops(opslice)

		err := l.Log(ops)
		if err != nil {
			panic(err)
		}
	}
	return logger
}

type logOpType uint64

const (
	createType uint64 = iota
	updateType
	deleteType
	connectType
	disconnectType
)

func opType(o Operation) uint64 {
	switch o.(type) {
	case *CreateOp:
		return createType
	case *UpdateOp:
		return updateType
	case *DeleteOp:
		return deleteType
	case *ConnectOp:
		return connectType
	case *DisconnectOp:
		return disconnectType
	default:
		panic("invalid optype")
	}
}

type Ops []Operation

func (o *Ops) Serialize() ([]byte, error) {
	w := newWriter()
	slice := []Operation(*o)
	sliceLen := uint64(len(slice))
	w.WriteUVarInt(sliceLen)

	for _, op := range slice {
		w.WriteUVarInt(opType(op))

		bytes, err := op.Serialize()
		if err != nil {
			return nil, err
		}

		w.WriteBytes(bytes)
	}
	return w.Done()
}

func (o *Ops) Deserialize(b *Builder, in []byte) (err error) {
	r := newReader(b, in)
	len := r.ReadUVarInt()

	for i := uint64(0); i < len; i++ {
		if r.Done() != nil {
			return r.Done()
		}
		opType := r.ReadUVarInt()
		opBytes := r.ReadBytes()

		switch opType {
		case createType:
			c := new(CreateOp)
			err = c.Deserialize(b, opBytes)
			*o = append(*o, c)
		case updateType:
			u := new(UpdateOp)
			err = u.Deserialize(b, opBytes)
			*o = append(*o, u)
		case deleteType:
			d := new(DeleteOp)
			err = d.Deserialize(b, opBytes)
			*o = append(*o, d)
		case connectType:
			c := new(ConnectOp)
			err = c.Deserialize(b, opBytes)
			*o = append(*o, c)
		case disconnectType:
			d := new(DisconnectOp)
			err = d.Deserialize(b, opBytes)
			*o = append(*o, d)
		default:
			panic("invalid optype")
		}
		if err != nil {
			return err
		}
	}

	return r.Done()
}

// conveience writer for avoiding lots of err == nil checks
type opWriter struct {
	err     error
	buf     bytes.Buffer
	scratch []byte
}

func newWriter() *opWriter {
	scratch := make([]byte, binary.MaxVarintLen64)
	return &opWriter{scratch: scratch}
}

func (l *opWriter) WriteID(id ID) {
	if l.err != nil {
		return
	}
	idBytes := id.Bytes()
	n, err := l.buf.Write(idBytes)
	if err != nil {
		l.err = err
		return
	}
	if n != 16 {
		l.err = errors.New("Bad write")
	}
	return
}

func (l *opWriter) WriteRecord(rec Record) {
	if l.err != nil {
		return
	}
	l.WriteID(rec.InterfaceID())
	l.WriteUVarInt(rec.Version())
	bytes, err := rec.MarshalBinary()
	if err != nil {
		l.err = err
	}
	l.WriteBytes(bytes)
	return
}

func (l *opWriter) WriteUVarInt(v uint64) {
	if l.err != nil {
		return
	}

	n := binary.PutUvarint(l.scratch, v)
	l.err = binary.Write(&l.buf, binary.LittleEndian, l.scratch[:n])
	return
}

func (l *opWriter) WriteBytes(bytes []byte) {
	if l.err != nil {
		return
	}

	bytelen := uint64(len(bytes))
	l.WriteUVarInt(bytelen)
	l.err = binary.Write(&l.buf, binary.LittleEndian, bytes)
	return
}

func (l *opWriter) Done() ([]byte, error) {
	return l.buf.Bytes(), l.err
}

func (l *opWriter) Err() error {
	return l.err
}

func newReader(b *Builder, data []byte) *opReader {
	return &opReader{buf: bytes.NewBuffer(data), builder: b}
}

type opReader struct {
	err     error
	buf     *bytes.Buffer
	builder *Builder
}

func (l *opReader) ReadID() (id ID) {
	if l.err != nil {
		return
	}
	bytes := make([]byte, 16)
	n, err := l.buf.Read(bytes)
	if err != nil {
		l.err = err
		return
	}
	if n != 16 {
		l.err = fmt.Errorf("read %v expected 16", n)
	}
	u, err := uuid.FromBytes(bytes)
	if err != nil {
		l.err = err
		return
	}
	id = ID(u)
	return
}

func (l *opReader) ReadRecord() (rec Record) {
	if l.err != nil {
		return
	}
	modelID := l.ReadID()
	version := l.ReadUVarInt()
	recData := l.ReadBytes()

	rec, err := l.builder.RecordForInterfaceVersion(modelID, version)
	if err != nil {
		l.err = err
		return
	}

	err = rec.UnmarshalBinary(recData)
	if err != nil {
		l.err = err
	}
	return
}

func (l *opReader) ReadUVarInt() (val uint64) {
	if l.err != nil {
		return
	}
	val, err := binary.ReadUvarint(l.buf)
	if err != nil {
		l.err = err
	}
	return
}

func (l *opReader) ReadBytes() (bytes []byte) {
	if l.err != nil {
		return
	}

	bytelen := l.ReadUVarInt()
	bytes = make([]byte, bytelen)
	l.err = binary.Read(l.buf, binary.LittleEndian, &bytes)
	return
}

func (l *opReader) Done() error {
	return l.err
}
