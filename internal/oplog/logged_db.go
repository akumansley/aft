package oplog

import (
	"fmt"

	"awans.org/aft/internal/db"
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

type TxEntry struct {
	Ops []DBOpEntry
}

func (txe TxEntry) String() string {
	return fmt.Sprintf("entry{%v}", txe.Ops)
}

func (txe TxEntry) isEmpty() bool {
	return len(txe.Ops) == 0
}

func (txe TxEntry) Replay(rwtx db.RWTx) {
	for _, op := range txe.Ops {
		op.Replay(rwtx)
	}
}

type DBOpEntry struct {
	OpType DBOp
	Op     interface{}
}

func (oe DBOpEntry) String() string {
	return fmt.Sprintf("op{%v}", oe.Op)
}

func (oe DBOpEntry) Replay(rwtx db.RWTx) {
	switch oe.OpType {
	case Create:
		cro := oe.Op.(CreateOp)
		cro.Replay(rwtx)
	case Connect:
		cno := oe.Op.(ConnectOp)
		cno.Replay(rwtx)
	case Update:
		uo := oe.Op.(UpdateOp)
		uo.Replay(rwtx)
	case Delete:
		do := oe.Op.(DeleteOp)
		do.Replay(rwtx)
	}
}

type CreateOp struct {
	RecFields interface{}
	ModelID   db.ID
}

func (op CreateOp) String() string {
	return fmt.Sprintf("create{%v %v}", op.RecFields, op.ModelID)
}

func (cro CreateOp) Replay(rwtx db.RWTx) {
	st := cro.RecFields
	m, err := rwtx.Schema().GetModelByID(cro.ModelID)
	if err != nil {
		panic("Model not found on replay")
	}
	rec := db.RecordFromParts(st, m)
	rwtx.Insert(rec)
}

type ConnectOp struct {
	Left  db.ID
	Right db.ID
	RelID db.ID
}

func (op ConnectOp) String() string {
	return fmt.Sprintf("connect{%v %v %v}", op.Left, op.Right, op.RelID)
}

func (cno ConnectOp) Replay(rwtx db.RWTx) {
	rwtx.Connect(cno.Left, cno.Right, cno.RelID)
}

type DisconnectOp struct {
	Left  db.ID
	Right db.ID
	RelID db.ID
}

func (op DisconnectOp) String() string {
	return fmt.Sprintf("disconnect{%v %v %v}", op.Left, op.Right, op.RelID)
}

func (dno DisconnectOp) Replay(rwtx db.RWTx) {
	rwtx.Disconnect(dno.Left, dno.Right, dno.RelID)
}

type UpdateOp struct {
	OldRecFields interface{}
	NewRecFields interface{}
	ModelID      db.ID
}

func (op UpdateOp) String() string {
	return fmt.Sprintf("update{%v %v %v}", op.OldRecFields, op.NewRecFields, op.ModelID)
}

func (uo UpdateOp) Replay(rwtx db.RWTx) {
	Ost := uo.OldRecFields
	Nst := uo.NewRecFields
	m, err := rwtx.Schema().GetModelByID(uo.ModelID)
	if err != nil {
		panic("Model not found on replay")
	}
	oldRec := db.RecordFromParts(Ost, m)
	newRec := db.RecordFromParts(Nst, m)
	rwtx.Update(oldRec, newRec)
}

type DeleteOp struct {
	RecFields interface{}
	ModelID   db.ID
}

func (op DeleteOp) String() string {
	return fmt.Sprintf("delete{%v %v}", op.RecFields, op.ModelID)
}

func (cro DeleteOp) Replay(rwtx db.RWTx) {
	st := cro.RecFields
	m, err := rwtx.Schema().GetModelByID(cro.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Delete(db.RecordFromParts(st, m))
}

type loggedDB struct {
	db.DB
	l OpLog
}

type loggedTx struct {
	db.RWTx
	txe TxEntry
	l   OpLog
}

func DBFromLog(db db.DB, l OpLog) error {
	iter := l.Iterator()
	for iter.Next() {
		rwtx := db.NewRWTx()
		val := iter.Value()
		txe := val.(TxEntry)
		txe.Replay(rwtx)
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

func LoggedDB(l OpLog, d db.DB) db.DB {
	return &loggedDB{DB: d, l: l}
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{RWTx: l.DB.NewRWTx(), l: l.l}
}

func (tx *loggedTx) Schema() *db.Schema {
	s := tx.RWTx.Schema()
	s.SetTx(tx)
	return s
}

func (tx *loggedTx) Insert(rec db.Record) error {
	co := CreateOp{RecFields: rec.RawData(), ModelID: rec.Interface().ID()}
	dboe := DBOpEntry{Create, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.RWTx.Insert(rec)
}

func (tx *loggedTx) Connect(left, right, rel db.ID) error {
	co := ConnectOp{Left: left, Right: right, RelID: rel}
	dboe := DBOpEntry{Connect, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.RWTx.Connect(left, right, rel)
}

func (tx *loggedTx) Disconnect(left, right, rel db.ID) error {
	co := DisconnectOp{Left: left, Right: right, RelID: rel}
	dboe := DBOpEntry{Disconnect, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.RWTx.Disconnect(left, right, rel)
}

func (tx *loggedTx) Update(oldRec, newRec db.Record) error {
	uo := UpdateOp{OldRecFields: oldRec.RawData(), NewRecFields: newRec.RawData(), ModelID: oldRec.Interface().ID()}
	dboe := DBOpEntry{Update, uo}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.RWTx.Update(oldRec, newRec)
}

func (tx *loggedTx) Delete(rec db.Record) error {
	co := DeleteOp{RecFields: rec.RawData(), ModelID: rec.Interface().ID()}
	dboe := DBOpEntry{Delete, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.RWTx.Delete(rec)
}

func (tx *loggedTx) Commit() (err error) {
	if !tx.txe.isEmpty() {
		err = tx.l.Log(tx.txe)
		if err != nil {
			return
		}
	}
	err = tx.RWTx.Commit()
	return
}
