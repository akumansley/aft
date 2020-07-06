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
	ModelID      db.ModelID
}

func (cro CreateOp) Replay(rwtx db.RWTx) {
	st := cro.RecFields
	m, err := rwtx.GetModelByID(cro.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Insert(db.RecordFromParts(st, m))
}

type ConnectOp struct {
	Left  db.ID
	Right db.ID
	RelID db.ID
}

func (cno ConnectOp) Replay(rwtx db.RWTx) {
	relRec, err := rwtx.FindOne(db.RelationshipModel.ID, db.EqID(cno.RelID))
	if err != nil {
		panic("couldn't find one on replay")
	}
	rel, err := db.LoadRel(relRec)
	if err != nil {
		panic("couldn't find one on replay")
	}
	left, err := rwtx.FindOne(rel.LeftModelID, db.EqID(cno.Left))
	if err != nil {
		panic("couldn't find one on replay")
	}
	right, err := rwtx.FindOne(rel.RightModelID, db.EqID(cno.Right))
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Connect(left, right, rel)
}

type UpdateOp struct {
	OldRecFields interface{}
	NewRecFields interface{}
	ModelID      db.ModelID
}

func (uo UpdateOp) Replay(rwtx db.RWTx) {
	Ost := uo.OldRecFields
	Nst := uo.NewRecFields
	m, err := rwtx.GetModelByID(uo.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Update(db.RecordFromParts(Ost, m), db.RecordFromParts(Nst, m))
}

type DeleteOp struct {
	RecFields interface{}
	ModelID   db.ModelID
}

func (cro DeleteOp) Replay(rwtx db.RWTx) {
	st := cro.RecFields
	m, err := rwtx.GetModelByID(cro.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Delete(db.RecordFromParts(st, m))
}

type loggedDB struct {
	inner db.DB
	l     OpLog
}

type loggedTx struct {
	inner db.RWTx
	txe   TxEntry
	l     OpLog
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
	return &loggedDB{inner: d, l: l}
}

func (l *loggedDB) NewTx() db.Tx {
	return l.inner.NewTx()
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{inner: l.inner.NewRWTx(), l: l.l}
}

func (l *loggedDB) DeepEquals(o db.DB) bool {
	return l.inner.DeepEquals(o)
}

func (l *loggedDB) Iterator() db.Iterator {
	return l.inner.Iterator()
}

func (tx *loggedTx) GetModelByID(id db.ModelID) (db.Model, error) {
	return tx.inner.GetModelByID(id)
}

func (tx *loggedTx) Ex() db.CodeExecutor {
	return tx.inner.Ex()
}

func (tx *loggedTx) GetModel(modelName string) (db.Model, error) {
	return tx.inner.GetModel(modelName)
}

func (tx *loggedTx) SaveModel(m db.Model) error {
	return tx.inner.SaveModel(m)
}

func (tx *loggedTx) Ref(u db.ModelID) db.ModelRef {
	return tx.inner.Ref(u)
}

func (tx *loggedTx) Query(m db.ModelRef) db.Q {
	return tx.inner.Query(m)
}

func (tx *loggedTx) MakeRecord(modelID db.ModelID) (db.Record, error) {
	return tx.inner.MakeRecord(modelID)
}

func (tx *loggedTx) Insert(rec db.Record) error {
	co := CreateOp{RecFields: rec.RawData(), ModelID: rec.Model().ID}
	dboe := DBOpEntry{Create, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Insert(rec)
}

func (tx *loggedTx) Connect(left, right db.Record, rel db.Relationship) error {
	co := ConnectOp{Left: left.ID(), Right: right.ID(), RelID: rel.ID}
	dboe := DBOpEntry{Connect, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Connect(left, right, rel)
}

func (tx *loggedTx) Update(oldRec, newRec db.Record) error {
	uo := UpdateOp{OldRecFields: oldRec.RawData(), NewRecFields: newRec.RawData(), ModelID: oldRec.Model().ID}
	dboe := DBOpEntry{Update, uo}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Update(oldRec, newRec)
}

func (tx *loggedTx) Delete(rec db.Record) error {
	co := DeleteOp{RecFields: rec.RawData(), ModelID: rec.Model().ID}
	dboe := DBOpEntry{Delete, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Delete(rec)
}

func (tx *loggedTx) FindOne(modelID db.ModelID, matcher db.Matcher) (db.Record, error) {
	return tx.inner.FindOne(modelID, matcher)
}

func (tx *loggedTx) FindMany(modelID db.ModelID, matcher db.Matcher) ([]db.Record, error) {
	return tx.inner.FindMany(modelID, matcher)
}

func (tx *loggedTx) Commit() (err error) {
	err = tx.l.Log(tx.txe)
	if err != nil {
		return
	}
	err = tx.inner.Commit()
	return
}
