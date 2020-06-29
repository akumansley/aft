package oplog

import (
	"awans.org/aft/internal/db"
)

type DBOp int

const (
	Connect DBOp = iota
	Create
	Update
	Delete
)

type TxEntry struct {
	Ops []DBOpEntry
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
	RecFields    interface{}
	RecordFields interface{}
	ModelID      db.ID
}

func (cro CreateOp) Replay(rwtx db.RWTx) {
	st := cro.RecFields
	m, err := rwtx.Schema().GetModelByID(cro.ModelID)
	if err != nil {
		panic("Model not found on replay")
	}
	rwtx.Insert(db.RecordFromParts(st, m))
}

type ConnectOp struct {
	Left  db.ID
	Right db.ID
	RelID db.ID
}

func (cno ConnectOp) Replay(rwtx db.RWTx) {
	rwtx.Connect(cno.Left, cno.Right, cno.RelID)
}

type UpdateOp struct {
	OldRecFields interface{}
	NewRecFields interface{}
	ModelID      db.ID
}

func (uo UpdateOp) Replay(rwtx db.RWTx) {
	Ost := uo.OldRecFields
	Nst := uo.NewRecFields
	m, err := rwtx.Schema().GetModelByID(uo.ModelID)
	if err != nil {
		panic("Model not found on replay")
	}
	rwtx.Update(db.RecordFromParts(Ost, m), db.RecordFromParts(Nst, m))
}

type DeleteOp struct {
	RecFields interface{}
	ModelID   db.ID
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
	rwtx := db.NewRWTx()
	for iter.Next() {
		val := iter.Value()
		txe := val.(TxEntry)
		txe.Replay(rwtx)
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	err := rwtx.Commit()
	return err
}

func LoggedDB(l OpLog, d db.DB) db.DB {
	return &loggedDB{DB: d, l: l}
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{RWTx: l.DB.NewRWTx(), l: l.l}
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
	err = tx.l.Log(tx.txe)
	if err != nil {
		return
	}
	err = tx.RWTx.Commit()
	return
}
