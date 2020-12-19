package db

import (
	"fmt"
)

type Operation interface {
	Replay(RWTx)
}

type CreateOp struct {
	Record  Record
	ModelID ID
}

func (op CreateOp) String() string {
	return fmt.Sprintf("create{%v %v}", op.Record.Map(), op.ModelID)
}

func (cro CreateOp) Replay(rwtx RWTx) {
	rwtx.Insert(cro.Record)
}

type ConnectOp struct {
	Source ID
	Target ID
	RelID  ID
}

func (op ConnectOp) String() string {
	return fmt.Sprintf("connect{%v %v %v}", op.Source, op.Target, op.RelID)
}

func (cno ConnectOp) Replay(rwtx RWTx) {
	rwtx.Connect(cno.Source, cno.Target, cno.RelID)
}

type DisconnectOp struct {
	Source ID
	Target ID
	RelID  ID
}

func (op DisconnectOp) String() string {
	return fmt.Sprintf("disconnect{%v %v %v}", op.Source, op.Target, op.RelID)
}

func (dno DisconnectOp) Replay(rwtx RWTx) {
	rwtx.Disconnect(dno.Source, dno.Target, dno.RelID)
}

type UpdateOp struct {
	OldRecord Record
	NewRecord Record
	ModelID   ID
}

func (op UpdateOp) String() string {
	return fmt.Sprintf("update{%v %v %v}", op.OldRecord.Map(), op.NewRecord.Map(), op.ModelID)
}

func (op UpdateOp) Replay(rwtx RWTx) {
	rwtx.Update(op.OldRecord, op.NewRecord)
}

type DeleteOp struct {
	Record  Record
	ModelID ID
}

func (op DeleteOp) String() string {
	return fmt.Sprintf("delete{%v %v}", op.Record.Map(), op.ModelID)
}

func (op DeleteOp) Replay(rwtx RWTx) {
	rwtx.Delete(op.Record)
}
