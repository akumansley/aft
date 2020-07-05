package db

import (
	"errors"
	"fmt"
)

var (
	ErrData         = errors.New("data-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

type Tx interface {
	Schema() *Schema
	loadFunction(Record) (Function, error)

	getRelatedOne(ID, ID, ID) (Record, error)
	getRelatedMany(ID, ID, ID) ([]Record, error)
	getRelatedManyReverse(ID, ID, ID) ([]Record, error)
	getRelatedOneReverse(ID, ID, ID) (Record, error)

	FindOne(ID, Matcher) (Record, error)
	FindMany(ID, Matcher) ([]Record, error)
	Ref(ID) ModelRef
	Query(ModelRef) Q
}

type RWTx interface {
	Schema() *Schema

	// reads
	getRelatedOne(ID, ID, ID) (Record, error)
	getRelatedMany(ID, ID, ID) ([]Record, error)
	getRelatedManyReverse(ID, ID, ID) ([]Record, error)
	getRelatedOneReverse(ID, ID, ID) (Record, error)

	FindOne(ID, Matcher) (Record, error)
	FindMany(ID, Matcher) ([]Record, error)
	Ref(ID) ModelRef
	Query(ModelRef) Q

	// writes
	MakeRecord(ID) (Record, error)
	Insert(Record) error
	Update(oldRec, newRec Record) error
	Delete(Record) error
	Connect(source, target, rel ID) error

	Commit() error
}

type holdTx struct {
	h     *Hold
	db    *holdDB
	rw    bool
	cache map[ID]interface{}
}

func (tx *holdTx) ensureWrite() {
	if !tx.rw {
		panic("Tried to write in a read only tx")
	}
}

func (tx *holdTx) loadFunction(rec Record) (f Function, err error) {
	mid := rec.Model().ID()
	rt := tx.db.runtimes[mid]
	f = rt.Load(tx, rec)
	return
}

func (tx *holdTx) FindOne(modelID ID, matcher Matcher) (rec Record, err error) {
	rec, err = tx.h.FindOne(modelID, matcher)
	return
}

func (tx *holdTx) FindMany(modelID ID, matcher Matcher) (recs []Record, err error) {
	recs, err = tx.h.FindMany(modelID, matcher)
	return
}

func (tx *holdTx) getRelatedOne(id, rel, target ID) (Record, error) {
	return tx.h.GetLinkedOne(id, rel, target)
}

func (tx *holdTx) getRelatedMany(id, rel, target ID) ([]Record, error) {
	return tx.h.GetLinkedMany(id, rel, target)
}

func (tx *holdTx) getRelatedManyReverse(id, rel, target ID) ([]Record, error) {
	return tx.h.GetLinkedManyReverse(id, rel, target)
}

func (tx *holdTx) getRelatedOneReverse(id ID, rel, target ID) (Record, error) {
	return tx.h.GetLinkedOneReverse(id, rel, target)
}

func (tx *holdTx) Insert(rec Record) error {
	tx.ensureWrite()
	tx.h = tx.h.Insert(rec)
	return nil
}

func (tx *holdTx) Update(oldRec, newRec Record) error {
	tx.ensureWrite()
	if oldRec.ID() != newRec.ID() {
		return fmt.Errorf("Can't update ID field on a record")
	}
	tx.h = tx.h.Insert(newRec)
	return nil
}

func (tx *holdTx) Connect(source, target, rel ID) error {
	tx.ensureWrite()
	// maybe unlink an existing relationship
	tx.h = tx.h.Link(source, target, rel)
	return nil
}

func (tx *holdTx) Delete(rec Record) error {
	tx.ensureWrite()
	tx.h = tx.h.Delete(rec)
	// todo: delete links
	return nil
}

func (tx *holdTx) MakeRecord(modelID ID) (rec Record, err error) {
	m := tx.Schema().GetModelByID(modelID)
	rec = RecordForModel(m)
	return
}

func (tx *holdTx) Schema() *Schema {
	return &Schema{tx}
}

func (tx *holdTx) Commit() error {
	tx.ensureWrite()
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()
	return nil
}
