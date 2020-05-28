package db

import (
	"awans.org/aft/internal/datatypes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
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

func (db *holdDB) AddMetaModel() {
	tx := db.NewRWTx()
	tx.SaveModel(ModelModel)
	tx.SaveModel(AttributeModel)
	tx.SaveModel(RelationshipModel)
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
}

type Tx interface {
	GetModel(string) (Model, error)
	GetModelById(uuid.UUID) (Model, error)
	MakeRecord(string) Record
	FindOne(string, Matcher) (Record, error)
	FindMany(string, Matcher) []Record
}

type RWTx interface {
	// remove
	GetModel(string) (Model, error)
	GetModelById(uuid.UUID) (Model, error)
	SaveModel(Model)

	FindOne(string, Matcher) (Record, error)
	FindMany(string, Matcher) []Record
	MakeRecord(string) Record

	// these are good, i think
	Insert(Record)
	Connect(from, to Record, fromRel Relationship)

	Commit() error
}

type holdDB struct {
	sync.RWMutex
	h *Hold
}

type holdTx struct {
	h  *Hold
	db *holdDB
	rw bool
}

func (tx *holdTx) ensureWrite() {
	if !tx.rw {
		panic("Tried to write in a read only tx")
	}
}

func (db *holdDB) NewTx() Tx {
	db.RLock()
	tx := holdTx{h: db.h, rw: false}
	db.RUnlock()
	return &tx
}

func (db *holdDB) NewRWTx() RWTx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: true}
	db.RUnlock()
	return &tx
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

func (tx *holdTx) FindOne(modelName string, matcher Matcher) (rec Record, err error) {
	rec, err = tx.h.FindOne(modelName, matcher)
	return
}

func (tx holdTx) FindMany(modelName string, matcher Matcher) []Record {
	mi := tx.h.IterMatches(modelName, matcher)
	hits := []Record{}
	for val, ok := mi.Next(); ok; val, ok = mi.Next() {
		hits = append(hits, val)
	}
	return hits
}

func (tx *holdTx) Insert(rec Record) {
	tx.ensureWrite()
	tx.h = tx.h.Insert(rec)
}

func (tx *holdTx) Connect(left, right Record, rel Relationship) {
	tx.ensureWrite()

	if rel.LeftBinding == BelongsTo && (rel.RightBinding == HasOne || rel.RightBinding == HasMany) {
		// FK left
		left.SetFK(rel.LeftName, right.Id())
	} else if rel.RightBinding == BelongsTo && (rel.LeftBinding == HasOne || rel.LeftBinding == HasMany) {
		// FK right
		right.SetFK(rel.RightName, left.Id())
	} else if rel.LeftBinding == HasManyAndBelongsToMany && rel.RightBinding == HasManyAndBelongsToMany {
		// Join table
		panic("Many to many relationships not implemented yet")
	} else {
		panic("Trying to connect invalid relationship")
	}
	h1 := tx.h.Insert(left)
	h2 := h1.Insert(right)
	tx.h = h2
}

func LoadRel(storeRel Record) Relationship {
	return Relationship{
		Id:           storeRel.Id(),
		LeftBinding:  RelType(storeRel.Get("leftBinding").(int64)),
		LeftModelId:  storeRel.GetFK("leftModel"),
		LeftName:     storeRel.Get("leftName").(string),
		RightBinding: RelType(storeRel.Get("rightBinding").(int64)),
		RightModelId: storeRel.GetFK("rightModel"),
		RightName:    storeRel.Get("rightName").(string),
	}
}

func loadModel(tx *holdTx, storeModel Record) Model {
	m := Model{
		Id:   storeModel.Id(),
		Name: storeModel.Get("name").(string),
	}

	attrs := make(map[string]Attribute)

	// make ModelId a dynamic key
	ami := tx.h.IterMatches("attribute", EqFK("model", m.Id))
	for storeAttr, ok := ami.Next(); ok; storeAttr, ok = ami.Next() {
		attr := Attribute{
			AttrType: datatypes.AttrType(storeAttr.Get("attrType").(int64)),
			Id:       storeAttr.Id(),
		}
		name := storeAttr.Get("name").(string)
		attrs[name] = attr
	}
	m.Attributes = attrs

	lRels := []Relationship{}
	rmi := tx.h.IterMatches("relationship", EqFK("leftModel", m.Id))
	for storeRel, ok := rmi.Next(); ok; storeRel, ok = rmi.Next() {
		lRels = append(lRels, LoadRel(storeRel))
	}
	m.LeftRelationships = lRels

	rRels := []Relationship{}
	rmi = tx.h.IterMatches("relationship", EqFK("rightModel", m.Id))
	for storeRel, ok := rmi.Next(); ok; storeRel, ok = rmi.Next() {
		rRels = append(rRels, LoadRel(storeRel))
	}
	m.RightRelationships = rRels
	return m
}

func (tx *holdTx) GetModelById(id uuid.UUID) (m Model, err error) {
	storeModel, err := tx.h.FindOne("model", Eq("id", id))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, id)
	}
	m = loadModel(tx, storeModel)

	return m, nil
}

func (tx *holdTx) GetModel(modelName string) (m Model, err error) {
	modelName = strings.ToLower(modelName)
	storeModel, err := tx.h.FindOne("model", Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	m = loadModel(tx, storeModel)

	return m, nil
}

func saveRel(tx *holdTx, rel Relationship) {
	storeRel := RecordForModel(RelationshipModel)
	storeRel.SetFK("leftModel", rel.LeftModelId)
	storeRel.Set("leftName", rel.LeftName)
	storeRel.Set("leftBinding", int64(rel.LeftBinding))
	storeRel.SetFK("rightModel", rel.RightModelId)
	storeRel.Set("rightName", rel.RightName)
	storeRel.Set("rightBinding", int64(rel.RightBinding))
	storeRel.Set("id", rel.Id)
	tx.h = tx.h.Insert(storeRel)
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m Model) {
	tx.ensureWrite()
	storeModel := RecordForModel(ModelModel)
	storeModel.Set("name", m.Name)
	storeModel.Set("id", m.Id)
	tx.h = tx.h.Insert(storeModel)

	for aKey, attr := range m.Attributes {
		storeAttr := RecordForModel(AttributeModel)
		storeAttr.Set("name", aKey)
		storeAttr.Set("attrType", int64(attr.AttrType))
		storeAttr.Set("id", attr.Id)
		storeAttr.SetFK("model", m.Id)
		tx.h = tx.h.Insert(storeAttr)
	}

	for _, rel := range m.RightRelationships {
		saveRel(tx, rel)
	}
	for _, rel := range m.LeftRelationships {
		saveRel(tx, rel)
	}
	// done for a side effect
	tx.MakeRecord(m.Name)
}

func (tx *holdTx) MakeRecord(modelName string) Record {
	modelName = strings.ToLower(modelName)
	m, _ := tx.GetModel(modelName)
	rec := RecordForModel(m)
	return rec
}

func (tx *holdTx) Commit() error {
	tx.ensureWrite()
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()
	return nil
}
