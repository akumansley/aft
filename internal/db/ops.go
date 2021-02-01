package db

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
)

type Operation interface {
	Replay(RWTx)
	Serialize() ([]byte, error)
	Deserialize(*Builder, []byte) error
}

type CreateOp struct {
	Record Record
}

func (c *CreateOp) String() string {
	return fmt.Sprintf("create{%v}", c.Record)
}

func (c *CreateOp) MarshalJSON() (bytes []byte, err error) {
	data := map[string]interface{}{}
	data["opType"] = "create"
	data["Record"] = c.Record
	if err != nil {
		return
	}
	return json.Marshal(data)
}

func (c *CreateOp) Replay(rwtx RWTx) {
	rwtx.Insert(c.Record)
}

func (c *CreateOp) Serialize() ([]byte, error) {
	w := newWriter()
	w.WriteRecord(c.Record)
	return w.Done()
}

func (c *CreateOp) Deserialize(b *Builder, data []byte) error {
	r := newReader(b, data)
	c.Record = r.ReadRecord()
	return r.Done()
}

type ConnectOp struct {
	Source ID
	Target ID
	RelID  ID
}

func (c *ConnectOp) String() string {
	return fmt.Sprintf("connect{%v %v %v}", c.Source, c.Target, c.RelID)
}

func (c *ConnectOp) MarshalJSON() ([]byte, error) {
	data := structs.Map(c)
	data["opType"] = "connect"
	return json.Marshal(data)
}

func (c *ConnectOp) Replay(rwtx RWTx) {
	rwtx.Connect(c.Source, c.Target, c.RelID)
}

func (c *ConnectOp) Serialize() ([]byte, error) {
	w := newWriter()
	w.WriteID(c.Source)
	w.WriteID(c.Target)
	w.WriteID(c.RelID)
	return w.Done()
}

func (c *ConnectOp) Deserialize(b *Builder, data []byte) error {
	r := newReader(b, data)
	c.Source = r.ReadID()
	c.Target = r.ReadID()
	c.RelID = r.ReadID()
	return r.Done()
}

type DisconnectOp struct {
	Source ID
	Target ID
	RelID  ID
}

func (op *DisconnectOp) String() string {
	return fmt.Sprintf("disconnect{%v %v %v}", op.Source, op.Target, op.RelID)
}

func (d *DisconnectOp) MarshalJSON() ([]byte, error) {
	data := structs.Map(d)
	data["opType"] = "disconnect"
	return json.Marshal(data)
}

func (d *DisconnectOp) Replay(rwtx RWTx) {
	rwtx.Disconnect(d.Source, d.Target, d.RelID)

}
func (d *DisconnectOp) Serialize() ([]byte, error) {
	w := newWriter()
	w.WriteID(d.Source)
	w.WriteID(d.Target)
	w.WriteID(d.RelID)
	return w.Done()
}

func (d *DisconnectOp) Deserialize(b *Builder, data []byte) error {
	r := newReader(b, data)
	d.Source = r.ReadID()
	d.Target = r.ReadID()
	d.RelID = r.ReadID()
	return r.Done()
}

type UpdateOp struct {
	OldRecord Record
	NewRecord Record
}

func (op *UpdateOp) String() string {
	return fmt.Sprintf("update{%v %v}", op.OldRecord, op.NewRecord)
}

func (op *UpdateOp) Replay(rwtx RWTx) {
	rwtx.Update(op.OldRecord, op.NewRecord)
}

func (u *UpdateOp) MarshalJSON() (bytes []byte, err error) {
	data := map[string]interface{}{}
	data["opType"] = "update"
	data["OldRecord"] = u.OldRecord
	data["NewRecord"] = u.NewRecord
	if err != nil {
		return
	}
	return json.Marshal(data)
}

func (u *UpdateOp) Serialize() ([]byte, error) {
	w := newWriter()
	w.WriteRecord(u.OldRecord)
	w.WriteRecord(u.NewRecord)
	return w.Done()
}

func (u *UpdateOp) Deserialize(b *Builder, data []byte) error {
	r := newReader(b, data)
	u.OldRecord = r.ReadRecord()
	u.NewRecord = r.ReadRecord()
	return r.Done()
}

type DeleteOp struct {
	Record Record
}

func (op *DeleteOp) String() string {
	return fmt.Sprintf("delete{%v}", op.Record)
}

func (d *DeleteOp) MarshalJSON() (bytes []byte, err error) {
	data := map[string]interface{}{}
	data["opType"] = "delete"
	data["Record"] = d.Record
	if err != nil {
		return
	}
	return json.Marshal(data)
}

func (op *DeleteOp) Replay(rwtx RWTx) {
	rwtx.Delete(op.Record)
}

func (d *DeleteOp) Serialize() ([]byte, error) {
	w := newWriter()
	w.WriteRecord(d.Record)
	return w.Done()
}

func (d *DeleteOp) Deserialize(b *Builder, data []byte) error {
	r := newReader(b, data)
	d.Record = r.ReadRecord()
	return r.Done()
}
