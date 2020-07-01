package db

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrData         = errors.New("data-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

func New() DB {
	appDB := holdDB{h: NewHold()}
	appDB.AddMetaModel()
	return &appDB
}

//tests only rely on golang execution
func NewTest() DB {
	return New()
}

func (db *holdDB) AddMetaModel() {
	// first add the native runtime
	nr := NewNativeRuntime()
	db.RegisterRuntime(nr)

	funcs := []NativeFunctionL{
		boolValidator,
		intValidator,
		stringValidator,
		uuidValidator,
		floatValidator,
	}

	for _, f := range funcs {
		nr.Save(f)
	}

	tx := db.NewRWTx()

	core := []Literal{
		Bool,
		Int,
		String,
		UUID,
		Float,
	}

	for _, d := range core {
		db.addLiteral(d)
	}

	models := []Literal{
		ModelModel,
		ConcreteAttributeModel,
		RelationshipModel,
		CoreDatatypeModel,
		EnumValueModel,
	}

	relationships := []Literal{
		ModelAttributes,
		RelationshipSource,
		RelationshipTarget,
		AttributeDatatype,
		DatatypeValidator,
	}

	for _, m := range models {
		db.addLiteral(m)
	}
	for _, r := range relationships {
		db.addLiteral(r)
	}

	tx.Commit()
}

type Iterator interface {
	Next() (Record, bool)
}

// DB is a value
type DB interface {
	NewTx() Tx
	NewRWTx() RWTx
	DeepEquals(DB) bool
	Iterator() Iterator
	RegisterRuntime(Runtime)
}

type Tx interface {
	Schema() *Schema

	GetRelatedOne(ID, Relationship) (Record, error)
	GetRelatedMany(ID, Relationship) ([]Record, error)
	GetRelatedManyReverse(ID, Relationship) ([]Record, error)
	FindOne(ID, Matcher) (Record, error)
	FindMany(ID, Matcher) ([]Record, error)
	Ref(ID) ModelRef
	Query(ModelRef) Q
}

type RWTx interface {
	// remove
	Schema() *Schema

	// reads
	GetRelatedOne(ID, Relationship) (Record, error)
	GetRelatedMany(ID, Relationship) ([]Record, error)
	GetRelatedManyReverse(ID, Relationship) ([]Record, error)
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

type holdDB struct {
	sync.RWMutex
	h        *Hold
	runtimes map[ID]Runtime
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

func (db *holdDB) NewTx() Tx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: false, cache: make(map[ID]interface{})}
	db.RUnlock()
	return &tx
}

func (db *holdDB) NewRWTx() RWTx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: true, cache: make(map[ID]interface{})}
	db.RUnlock()
	return &tx
}

func (db *holdDB) RegisterRuntime(r Runtime) {
	m := r.ProvideModel()
	db.addLiteral(m)
	db.runtimes[m.ID] = r
	r.Registered(db)
}

func (db *holdDB) addLiteral(lit Literal) {
	tx := db.NewRWTx()
	recs, links := lit.MarshalDB()
	for _, rec := range recs {
		tx.Insert(rec)
	}
	for _, link := range links {
		tx.Connect(link.from, link.to, link.rel.ID)
	}
	tx.Commit()
}

func (db *holdDB) Iterator() Iterator {
	return db.h.Iterator()
}

func (db *holdDB) DeepEquals(o DB) bool {
	leftI := db.Iterator()
	rightI := o.Iterator()
	for {
		lR, lok := leftI.Next()
		rR, rok := rightI.Next()
		if lok != rok {
			return false
		}
		if lok {
			if !lR.DeepEquals(rR) {
				return false
			}
		} else {
			return true
		}
	}
}

func (tx *holdTx) FindOne(modelID ID, matcher Matcher) (rec Record, err error) {
	rec, err = tx.h.FindOne(modelID, matcher)
	return
}

func (tx *holdTx) FindMany(modelID ID, matcher Matcher) (recs []Record, err error) {
	recs, err = tx.h.FindMany(modelID, matcher)
	return
}

func (tx *holdTx) GetRelatedOne(id ID, rel Relationship) (Record, error) {
	return tx.h.GetLinkedOne(id, rel)
}

func (tx *holdTx) GetRelatedMany(id ID, rel Relationship) ([]Record, error) {
	return tx.h.GetLinkedMany(id, rel)
}

func (tx *holdTx) GetRelatedManyReverse(id ID, rel Relationship) ([]Record, error) {
	return tx.h.GetLinkedManyReverse(id, rel)
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
