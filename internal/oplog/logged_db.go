package oplog

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"awans.org/aft/internal/db"
	"github.com/fatih/structs"
	"github.com/google/uuid"
)

func DBFromLog(appDB db.DB, l OpLog) error {
	iter := l.Iterator()
	for iter.Next() {
		rwtx := appDB.NewRWTx()
		var logOps LogOps
		val := iter.Value()
		bts := val.([]byte)
		err := logOps.UnmarshalBinary(bts)
		if err != nil {
			return err
		}
		var ops []db.Operation
		for _, logOp := range logOps {
			op, err := logOp.Decode(rwtx)
			if err != nil {
				return err
			}
			ops = append(ops, op)
		}
		for _, op := range ops {
			op.Replay(rwtx)
		}
		err = rwtx.Commit()
		if err != nil {
			return err
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	return nil
}

func modelForAttr(tx db.Tx, attrID db.ID) db.Model {
	models := tx.Ref(db.ModelModel.ID())
	attrs := tx.Ref(db.ConcreteAttributeModel.ID())
	q := tx.Query(models,
		db.Join(attrs, models.Rel(db.ModelAttributes)),
		db.Aggregate(attrs, db.Some),
		db.Filter(attrs, db.EqID(attrID)))
	rec, err := q.OneRecord()
	if err != nil {
		panic(err)
	}
	model := tx.Schema().LoadModel(rec)
	return model
}

func MakeTransactionLogger(l OpLog) func(db.BeforeCommit) {
	logger := func(event db.BeforeCommit) {
		ops := event.Tx.Operations()

		if len(ops) == 0 {
			return
		}
		encoded, err := encode(ops)
		if err != nil {
			panic(err)
		}

		err = l.Log(encoded)
		if err != nil {
			panic(err)
		}
	}
	return logger
}

type logOpType int16

const (
	createType logOpType = iota
	updateType
	deleteType
	connectType
	disconnectType
)

type LogOps []LogOp

func (o *LogOps) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	slice := []LogOp(*o)
	sliceLen := int64(len(slice))

	err := binary.Write(&buf, binary.LittleEndian, sliceLen)
	for _, op := range slice {
		err := binary.Write(&buf, binary.LittleEndian, op.Type())
		if err != nil {
			return nil, err
		}

		bytes, err := op.MarshalBinary()
		if err != nil {
			return nil, err
		}

		logOpLen := int64(len(bytes))
		err = binary.Write(&buf, binary.LittleEndian, logOpLen)
		err = binary.Write(&buf, binary.LittleEndian, bytes)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), err
}

func (o *LogOps) UnmarshalBinary(in []byte) error {
	var len int64
	buf := bytes.NewBuffer(in)
	err := binary.Read(buf, binary.LittleEndian, &len)
	if err != nil {
		return fmt.Errorf("log-error: err reading LogOps len %w", err)
	}

	for i := int64(0); i < len; i++ {
		var opTypeInt int16
		err = binary.Read(buf, binary.LittleEndian, &opTypeInt)
		if err != nil {
			return err
		}
		opType := logOpType(opTypeInt)

		var logOpLen int64
		err = binary.Read(buf, binary.LittleEndian, &logOpLen)
		if err != nil {
			return err
		}

		opBytes := make([]byte, logOpLen)
		err = binary.Read(buf, binary.LittleEndian, &opBytes)
		if err != nil {
			return err
		}

		switch opType {
		case createType:
			c := new(create)
			err = c.UnmarshalBinary(opBytes)
			*o = append(*o, c)
		case updateType:
			u := new(update)
			err = u.UnmarshalBinary(opBytes)
			*o = append(*o, u)
		case deleteType:
			d := new(update)
			err = d.UnmarshalBinary(opBytes)
			*o = append(*o, d)
		case connectType:
			c := new(connect)
			err = c.UnmarshalBinary(opBytes)
			*o = append(*o, c)
		case disconnectType:
			d := new(disconnect)
			err = d.UnmarshalBinary(opBytes)
			*o = append(*o, d)
		default:
			panic("invalid optype")
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func encode(ops []db.Operation) (*LogOps, error) {
	var encoded []LogOp

	for _, op := range ops {
		switch op.(type) {
		case db.CreateOp:
			c := op.(db.CreateOp)
			recBytes, err := c.Record.MarshalBinary()
			if err != nil {
				return nil, err
			}
			enc := &create{recBytes, c.ModelID}
			encoded = append(encoded, enc)
		case db.DeleteOp:
			d := op.(db.DeleteOp)
			recBytes, err := d.Record.MarshalBinary()
			if err != nil {
				return nil, err
			}
			enc := &delete{recBytes, d.ModelID}
			encoded = append(encoded, enc)
		case db.UpdateOp:
			u := op.(db.UpdateOp)
			oldRecBytes, err := u.OldRecord.MarshalBinary()
			if err != nil {
				return nil, err
			}
			newRecBytes, err := u.NewRecord.MarshalBinary()
			if err != nil {
				return nil, err
			}
			enc := &update{oldRecBytes, newRecBytes, u.ModelID}
			encoded = append(encoded, enc)
		case db.ConnectOp:
			c := op.(db.ConnectOp)
			enc := &connect{c.Source, c.Target, c.RelID}
			encoded = append(encoded, enc)
		case db.DisconnectOp:
			d := op.(db.DisconnectOp)
			enc := &disconnect{d.Source, d.Target, d.RelID}
			encoded = append(encoded, enc)
		default:
			panic(errors.New("invalid op"))
		}
	}
	logOps := LogOps(encoded)
	return &logOps, nil
}

type LogOp interface {
	Decode(db.Tx) (db.Operation, error)
	String() string
	MarshalBinary() ([]byte, error)
	Type() logOpType
}

// conveience writer for avoiding lots of err == nil checks
type loWriter struct {
	err error
	buf bytes.Buffer
}

func newLOWriter() *loWriter {
	return &loWriter{}
}

func (l *loWriter) WriteID(id db.ID) {
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

func (l *loWriter) Write(bytes []byte) {
	if l.err != nil {
		return
	}
	bytelen := int64(len(bytes))
	err := binary.Write(&l.buf, binary.LittleEndian, bytelen)
	if err != nil {
		l.err = err
		return
	}

	l.err = binary.Write(&l.buf, binary.LittleEndian, bytes)
	return
}

func (l *loWriter) Done() ([]byte, error) {
	return l.buf.Bytes(), l.err
}

func newLOReader(b []byte) *loReader {
	return &loReader{buf: bytes.NewBuffer(b)}
}

type loReader struct {
	err error
	buf *bytes.Buffer
}

func (l *loReader) ReadID() (id db.ID) {
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
	id = db.ID(u)
	return
}

func (l *loReader) ReadBytes() (bytes []byte) {
	if l.err != nil {
		return
	}
	var bytelen int64
	err := binary.Read(l.buf, binary.LittleEndian, &bytelen)
	if err != nil {
		l.err = err
		return
	}
	bytes = make([]byte, bytelen)
	l.err = binary.Read(l.buf, binary.LittleEndian, &bytes)
	return
}

func (l *loReader) Done() error {
	return l.err
}

type create struct {
	RecData []byte
	ModelID db.ID
}

func (c *create) Decode(tx db.Tx) (op db.Operation, err error) {
	model, err := tx.Schema().GetModelByID(c.ModelID)
	if err != nil {
		return nil, err
	}
	rec := db.RecordForModel(model)
	err = rec.UnmarshalBinary(c.RecData)
	if err != nil {
		return
	}
	op = db.CreateOp{
		Record:  rec,
		ModelID: c.ModelID,
	}
	return
}

func (c *create) MarshalBinary() ([]byte, error) {
	w := newLOWriter()
	w.WriteID(c.ModelID)
	w.Write(c.RecData)
	return w.Done()
}

func (c *create) UnmarshalBinary(bytes []byte) error {
	r := newLOReader(bytes)
	c.ModelID = r.ReadID()
	c.RecData = r.ReadBytes()
	return r.Done()
}

func (c *create) MarshalJSON() ([]byte, error) {
	data := structs.Map(c)
	data["opType"] = "create"
	return json.Marshal(data)
}

func (c *create) String() string {
	bytes, _ := c.MarshalJSON()
	return string(bytes)
}

func (c *create) Type() logOpType {
	return createType
}

type update struct {
	OldRecData []byte
	NewRecData []byte
	ModelID    db.ID
}

func (u *update) MarshalBinary() ([]byte, error) {
	w := newLOWriter()
	w.WriteID(u.ModelID)
	w.Write(u.OldRecData)
	w.Write(u.NewRecData)
	return w.Done()
}

func (u *update) UnmarshalBinary(bytes []byte) error {
	r := newLOReader(bytes)
	u.ModelID = r.ReadID()
	u.OldRecData = r.ReadBytes()
	u.NewRecData = r.ReadBytes()
	return r.Done()
}

func (u *update) MarshalJSON() ([]byte, error) {
	data := structs.Map(u)
	data["opType"] = "update"
	return json.Marshal(data)
}

func (u *update) String() string {
	bytes, _ := u.MarshalJSON()
	return string(bytes)
}

func (u *update) Decode(tx db.Tx) (op db.Operation, err error) {
	model, err := tx.Schema().GetModelByID(u.ModelID)
	if err != nil {
		return nil, err
	}
	oldRec := db.RecordForModel(model)
	err = oldRec.UnmarshalBinary(u.OldRecData)
	if err != nil {
		return
	}

	newRec := db.RecordForModel(model)
	err = newRec.UnmarshalBinary(u.NewRecData)
	if err != nil {
		return
	}
	return db.UpdateOp{
		OldRecord: oldRec,
		NewRecord: newRec,
		ModelID:   u.ModelID,
	}, nil
}

func (u *update) Type() logOpType {
	return updateType
}

type delete struct {
	RecData []byte
	ModelID db.ID
}

func (d *delete) UnmarshalBinary(bytes []byte) error {
	r := newLOReader(bytes)
	d.ModelID = r.ReadID()
	d.RecData = r.ReadBytes()
	return r.Done()
}

func (d *delete) MarshalBinary() ([]byte, error) {
	w := newLOWriter()
	w.WriteID(d.ModelID)
	w.Write(d.RecData)
	return w.Done()
}

func (d *delete) Decode(tx db.Tx) (op db.Operation, err error) {
	model, err := tx.Schema().GetModelByID(d.ModelID)
	if err != nil {
		return nil, err
	}
	rec := db.RecordForModel(model)
	err = rec.UnmarshalBinary(d.RecData)
	if err != nil {
		return
	}

	op = db.DeleteOp{
		Record:  db.RecordFromParts(d.RecData, model),
		ModelID: d.ModelID,
	}
	return
}
func (d *delete) MarshalJSON() ([]byte, error) {
	data := structs.Map(d)
	data["opType"] = "delete"
	return json.Marshal(data)
}
func (c *delete) String() string {
	bytes, _ := c.MarshalJSON()
	return string(bytes)
}

func (d *delete) Type() logOpType {
	return deleteType
}

type connect struct {
	Source db.ID
	Target db.ID
	RelID  db.ID
}

func (c *connect) MarshalBinary() ([]byte, error) {
	w := newLOWriter()
	w.WriteID(c.Source)
	w.WriteID(c.Target)
	w.WriteID(c.RelID)
	return w.Done()
}

func (c *connect) UnmarshalBinary(bytes []byte) error {
	r := newLOReader(bytes)
	c.Source = r.ReadID()
	c.Target = r.ReadID()
	c.RelID = r.ReadID()
	return r.Done()
}

func (c *connect) Decode(tx db.Tx) (db.Operation, error) {
	return db.ConnectOp{
		Source: c.Source, Target: c.Target, RelID: c.RelID,
	}, nil
}
func (c *connect) MarshalJSON() ([]byte, error) {
	data := structs.Map(c)
	data["opType"] = "connect"
	return json.Marshal(data)
}

func (c *connect) String() string {
	bytes, _ := c.MarshalJSON()
	return string(bytes)
}

func (c *connect) Type() logOpType {
	return connectType
}

type disconnect struct {
	Source db.ID
	Target db.ID
	RelID  db.ID
}

func (d *disconnect) MarshalBinary() ([]byte, error) {
	w := newLOWriter()
	w.WriteID(d.Source)
	w.WriteID(d.Target)
	w.WriteID(d.RelID)
	return w.Done()
}

func (d *disconnect) UnmarshalBinary(bytes []byte) error {
	r := newLOReader(bytes)
	d.Source = r.ReadID()
	d.Target = r.ReadID()
	d.RelID = r.ReadID()
	return r.Done()
}

func (d *disconnect) MarshalJSON() ([]byte, error) {
	data := structs.Map(d)
	data["opType"] = "disconnect"
	return json.Marshal(data)
}

func (d *disconnect) String() string {
	bytes, _ := d.MarshalJSON()
	return string(bytes)
}

func (d *disconnect) Decode(tx db.Tx) (db.Operation, error) {
	return db.DisconnectOp{
		Source: d.Source, Target: d.Target, RelID: d.RelID,
	}, nil
}

func (d *disconnect) Type() logOpType {
	return disconnectType
}
