package db

import (
	"fmt"
)

type DBOp int

const (
	Connect DBOp = iota
	Disconnect
	Create
	Update
	Delete
)

func (op DBOp) String() string {
	switch op {
	case Connect:
		return "connect"
	case Disconnect:
		return "disconnect"
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	default:
		panic("unknown op")
	}
}

type Operation interface {
	Replay(RWTx)
}

type CreateOp struct {
	RecFields interface{}
	ModelID   ID
}

func (op CreateOp) String() string {
	return fmt.Sprintf("create{%v %v}", op.RecFields, op.ModelID)
}

func (cro CreateOp) Replay(rwtx RWTx) {
	st := cro.RecFields
	m, err := rwtx.Schema().GetModelByID(cro.ModelID)
	if err != nil {
		panic("Model not found on replay")
	}
	rec := RecordFromParts(st, m)
	rwtx.Insert(rec)
}

type ConnectOp struct {
	Left  ID
	Right ID
	RelID ID
}

func (op ConnectOp) String() string {
	return fmt.Sprintf("connect{%v %v %v}", op.Left, op.Right, op.RelID)
}

func (cno ConnectOp) Replay(rwtx RWTx) {
	rwtx.Connect(cno.Left, cno.Right, cno.RelID)
}

type DisconnectOp struct {
	Left  ID
	Right ID
	RelID ID
}

func (op DisconnectOp) String() string {
	return fmt.Sprintf("disconnect{%v %v %v}", op.Left, op.Right, op.RelID)
}

func (dno DisconnectOp) Replay(rwtx RWTx) {
	rwtx.Disconnect(dno.Left, dno.Right, dno.RelID)
}

type UpdateOp struct {
	OldRecFields interface{}
	NewRecFields interface{}
	ModelID      ID
}

func (op UpdateOp) String() string {
	return fmt.Sprintf("update{%v %v %v}", op.OldRecFields, op.NewRecFields, op.ModelID)
}

func (uo UpdateOp) Replay(rwtx RWTx) {
	Ost := uo.OldRecFields
	Nst := uo.NewRecFields
	m, err := rwtx.Schema().GetModelByID(uo.ModelID)
	if err != nil {
		panic("Model not found on replay")
	}
	oldRec := RecordFromParts(Ost, m)
	newRec := RecordFromParts(Nst, m)
	rwtx.Update(oldRec, newRec)
}

type DeleteOp struct {
	RecFields interface{}
	ModelID   ID
}

func (op DeleteOp) String() string {
	return fmt.Sprintf("delete{%v %v}", op.RecFields, op.ModelID)
}

func (op DeleteOp) Replay(rwtx RWTx) {
	st := op.RecFields
	m, err := rwtx.Schema().GetModelByID(op.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Delete(RecordFromParts(st, m))
}
