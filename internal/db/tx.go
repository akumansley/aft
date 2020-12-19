package db

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrData         = errors.New("data-error")
	ErrConstraint   = errors.New("constraint-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

type Tx interface {
	AsOfStart() Tx
	Schema() *Schema
	Context() context.Context

	loadFunction(Record) (Function, error)
	getRelatedOne(ID, ID) (Record, error)
	getRelatedMany(ID, ID) ([]Record, error)
	getRelatedManyReverse(ID, ID) ([]Record, error)
	getRelatedOneReverse(ID, ID) (Record, error)

	MakeRecord(ID) (Record, error)
	Invalidate(ID) (Record, error)
	Ref(ID) ModelRef
	Query(ModelRef, ...QueryClause) Q
	Operations() []Operation

	Commit() error
	Abort(error)
}

type RWTx interface {
	Tx
	Insert(Record) error
	Update(oldRec, newRec Record) error
	Delete(Record) error
	Connect(source, target, rel ID) error
	Disconnect(source, target, rel ID) error

	unloggedUpdate(Record, Record) error
	unloggedDelete(Record) error
	unloggedDisconnect(ID, ID, ID) error
}

type holdTx struct {
	initH   *Hold
	h       *Hold
	db      *holdDB
	rw      bool
	ops     []Operation
	ctx     context.Context
	aborted error
}

func (tx *holdTx) Abort(err error) {
	tx.aborted = err
}

func (tx *holdTx) Context() context.Context {
	return tx.ctx
}

func (tx *holdTx) Operations() []Operation {
	return tx.ops
}

func (tx *holdTx) ensureWrite() {
	if !tx.rw {
		panic("Tried to write in a read only tx")
	}
}

func (tx *holdTx) AsOfStart() Tx {
	return &holdTx{initH: tx.initH, h: tx.initH, db: tx.db, rw: false, ops: nil, ctx: tx.ctx, aborted: nil}
}

func (tx *holdTx) loadFunction(rec Record) (f Function, err error) {
	mid := rec.Interface().ID()
	rt := tx.db.runtimes[mid]
	f = rt.Load(tx, rec)
	return
}

func (tx *holdTx) getRelatedOne(id, rel ID) (Record, error) {
	r, err := tx.h.GetLinkedOne(id, rel)
	return r, err
}

func (tx *holdTx) getRelatedMany(id, rel ID) ([]Record, error) {
	return tx.h.GetLinkedMany(id, rel)
}

func (tx *holdTx) getRelatedManyReverse(id, rel ID) ([]Record, error) {
	return tx.h.GetLinkedManyReverse(id, rel)
}

func (tx *holdTx) getRelatedOneReverse(id, rel ID) (Record, error) {
	return tx.h.GetLinkedOneReverse(id, rel)
}

func (tx *holdTx) Insert(rec Record) error {
	tx.ensureWrite()
	tx.h = tx.h.Insert(rec)
	co := CreateOp{Record: rec, ModelID: rec.Interface().ID()}
	tx.ops = append(tx.ops, co)
	return nil
}

func (tx *holdTx) Update(oldRec, newRec Record) error {
	err := tx.unloggedUpdate(oldRec, newRec)
	if err != nil {
		return err
	}
	uo := UpdateOp{OldRecord: oldRec, NewRecord: newRec, ModelID: oldRec.Interface().ID()}
	tx.ops = append(tx.ops, uo)
	return nil
}

func (tx *holdTx) unloggedUpdate(oldRec, newRec Record) error {
	tx.ensureWrite()
	if oldRec.ID() != newRec.ID() {
		return fmt.Errorf("Can't update ID field on a record")
	}
	tx.h = tx.h.Insert(newRec)
	return nil
}

func (tx *holdTx) Connect(source, target, relID ID) error {
	tx.ensureWrite()
	rel, err := tx.Schema().GetRelationshipByID(relID)
	if err != nil {
		return err
	}
	if !rel.Multi() {
		v, err := tx.getRelatedOne(source, relID)
		if err == nil {
			return fmt.Errorf("%w: can't connect already-connected (%v) non-multi relationship %v", ErrConstraint, v.ID(), rel.Name())
		}
		if !errors.Is(ErrNotFound, err) {
			return err
		}
	}
	tx.h = tx.h.Link(source, target, relID)

	co := ConnectOp{Source: source, Target: target, RelID: relID}
	tx.ops = append(tx.ops, co)
	return nil
}

func (tx *holdTx) Disconnect(source, target, relID ID) error {
	err := tx.unloggedDisconnect(source, target, relID)
	if err != nil {
		return err
	}
	do := DisconnectOp{Source: source, Target: target, RelID: relID}
	tx.ops = append(tx.ops, do)
	return nil
}

func (tx *holdTx) unloggedDisconnect(source, target, relID ID) error {
	tx.ensureWrite()
	tx.h = tx.h.Unlink(source, target, relID)
	return nil
}

func (tx *holdTx) unloggedDelete(rec Record) error {
	tx.ensureWrite()

	// all of this can get moved into the hold
	// and implemented with a scan of the rel indexes
	recIf := rec.Interface()
	cRel := tx.Ref(ConcreteRelationshipModel.ID())
	ifif := tx.Ref(InterfaceInterface.ID())
	target, _ := tx.Schema().GetRelationshipByID(ConcreteRelationshipTarget.ID())
	source, _ := tx.Schema().GetRelationshipByID(ConcreteRelationshipSource.ID())

	targetRels := tx.Query(cRel,
		Join(ifif, cRel.Rel(target)),
		Filter(ifif, EqID(recIf.ID())),
	).Records()

	sourceRels := tx.Query(cRel,
		Join(ifif, cRel.Rel(source)),
		Filter(ifif, EqID(recIf.ID())),
	).Records()

	rels := append(targetRels, sourceRels...)

	for _, relRec := range rels {
		rel, err := tx.Schema().loadRelationship(relRec)
		if err != nil {
			return err
		}
		if rel.Multi() {
			var targets []Record
			targets, _ = tx.getRelatedMany(rec.ID(), rel.ID())
			for _, tar := range targets {
				tx.Disconnect(rec.ID(), tar.ID(), rel.ID())
			}
			targets, _ = tx.getRelatedManyReverse(rec.ID(), rel.ID())
			for _, tar := range targets {
				tx.Disconnect(tar.ID(), rec.ID(), rel.ID())
			}
		} else {
			var target Record
			target, _ = tx.getRelatedOne(rec.ID(), rel.ID())
			if target != nil {
				tx.Disconnect(rec.ID(), target.ID(), rel.ID())
			}
			target, _ = tx.getRelatedOneReverse(rec.ID(), rel.ID())
			if target != nil {
				tx.Disconnect(target.ID(), rec.ID(), rel.ID())
			}
		}
	}
	tx.h = tx.h.Delete(rec)
	return nil
}

func (tx *holdTx) Delete(rec Record) error {
	err := tx.unloggedDelete(rec)
	if err != nil {
		return err
	}
	do := DeleteOp{Record: rec, ModelID: rec.Interface().ID()}
	tx.ops = append(tx.ops, do)
	return nil
}

func (tx *holdTx) Invalidate(interfaceID ID) (rec Record, err error) {
	i, err := tx.Schema().GetInterfaceByID(interfaceID)
	if err != nil {
		return
	}
	interfaceUpdated(i)
	return
}

func (tx *holdTx) MakeRecord(interfaceID ID) (rec Record, err error) {
	i, err := tx.Schema().GetInterfaceByID(interfaceID)
	if err != nil {
		return
	}
	rec = RecordForModel(i)
	return
}

func (tx *holdTx) Schema() *Schema {
	return &Schema{tx, tx.db}
}

func (tx *holdTx) Commit() error {
	tx.ensureWrite()
	if tx.aborted != nil {
		return tx.aborted
	}

	tx.db.bus.Publish(BeforeCommit{tx})
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()

	return nil
}

func (tx *holdTx) String() string {
	return tx.h.String()
}
