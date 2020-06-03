package db

import (
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
	//Add native datatypes and their code execution to the tree. Comes before models.
	for _, v := range datatypeMap {
		r := RecordForModel(DatatypeModel)
		SaveDatatype(r, v)
		tx.Insert(r)
	}
	for _, v := range codeMap {
		r := RecordForModel(CodeModel)
		SaveCode(r, v)
		tx.Insert(r)
	}

	tx.SaveModel(ModelModel)
	tx.SaveModel(AttributeModel)
	tx.SaveModel(RelationshipModel)
	tx.SaveModel(DatatypeModel)
	tx.SaveModel(CodeModel)

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
	GetModelByID(uuid.UUID) (Model, error)
	MakeRecord(string) Record
	FindOne(string, Matcher) (Record, error)
	FindMany(string, Matcher) []Record
}

type RWTx interface {
	// remove
	GetModel(string) (Model, error)
	GetModelByID(uuid.UUID) (Model, error)
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
		left.SetFK(rel.LeftName, right.ID())
	} else if rel.RightBinding == BelongsTo && (rel.LeftBinding == HasOne || rel.LeftBinding == HasMany) {
		// FK right
		right.SetFK(rel.RightName, left.ID())
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
		ID:           storeRel.ID(),
		LeftBinding:  RelType(storeRel.Get("leftBinding").(int64)),
		LeftModelID:  storeRel.GetFK("leftModel"),
		LeftName:     storeRel.Get("leftName").(string),
		RightBinding: RelType(storeRel.Get("rightBinding").(int64)),
		RightModelID: storeRel.GetFK("rightModel"),
		RightName:    storeRel.Get("rightName").(string),
	}
}

func loadModel(tx *holdTx, storeModel Record) Model {
	m := Model{
		ID:   storeModel.ID(),
		Name: storeModel.Get("name").(string),
	}

	attrs := make(map[string]Attribute)

	// make ModelID a dynamic key
	ami := tx.h.IterMatches("attribute", EqFK("model", m.ID))
	for storeAttr, ok := ami.Next(); ok; storeAttr, ok = ami.Next() {
		storeDatatype, _ := tx.h.FindOne("datatype", Eq("id", storeAttr.GetFK("datatype")))
		storeValidator, _ := tx.h.FindOne("code", Eq("id", storeDatatype.GetFK("validator")))
		validator := Code{
			ID:       storeValidator.ID(),
			Name:     storeValidator.Get("name").(string),
			Runtime:  Runtime(storeValidator.Get("runtime").(int64)),
			Code:     storeValidator.Get("code").(string),
			Function: Function(storeValidator.Get("function").(int64)),
		}
		d := Datatype{
			ID:          storeDatatype.ID(),
			Name:        storeDatatype.Get("name").(string),
			Validator:   validator,
			StorageType: StorageType(storeDatatype.Get("storageType").(int64)),
		}
		attr := Attribute{
			Datatype: d,
			ID:       storeAttr.ID(),
		}
		name := storeAttr.Get("name").(string)
		attrs[name] = attr
	}
	m.Attributes = attrs

	lRels := []Relationship{}
	rmi := tx.h.IterMatches("relationship", EqFK("leftModel", m.ID))
	for storeRel, ok := rmi.Next(); ok; storeRel, ok = rmi.Next() {
		lRels = append(lRels, LoadRel(storeRel))
	}
	m.LeftRelationships = lRels

	rRels := []Relationship{}
	rmi = tx.h.IterMatches("relationship", EqFK("rightModel", m.ID))
	for storeRel, ok := rmi.Next(); ok; storeRel, ok = rmi.Next() {
		rRels = append(rRels, LoadRel(storeRel))
	}
	m.RightRelationships = rRels
	return m
}

func (tx *holdTx) GetModelByID(id uuid.UUID) (m Model, err error) {
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

func SaveDatatype(storeDatatype Record, d Datatype) {
	storeDatatype.Set("id", d.ID)
	storeDatatype.Set("name", d.Name)
	storeDatatype.Set("storageType", int64(d.StorageType))
	storeDatatype.SetFK("validator", d.Validator.ID)
}

func SaveCode(storeCode Record, c Code) {
	storeCode.Set("id", c.ID)
	storeCode.Set("name", c.Name)
	storeCode.Set("runtime", int64(c.Runtime))
	storeCode.Set("function", int64(c.Function))
	storeCode.Set("code", c.Code)
}

func saveRel(tx *holdTx, rel Relationship) {
	storeRel := RecordForModel(RelationshipModel)
	storeRel.SetFK("leftModel", rel.LeftModelID)
	storeRel.Set("leftName", rel.LeftName)
	storeRel.Set("leftBinding", int64(rel.LeftBinding))
	storeRel.SetFK("rightModel", rel.RightModelID)
	storeRel.Set("rightName", rel.RightName)
	storeRel.Set("rightBinding", int64(rel.RightBinding))
	storeRel.Set("id", rel.ID)
	tx.h = tx.h.Insert(storeRel)
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m Model) {
	tx.ensureWrite()
	storeModel := RecordForModel(ModelModel)
	storeModel.Set("name", m.Name)
	storeModel.Set("id", m.ID)
	tx.h = tx.h.Insert(storeModel)

	for aKey, attr := range m.Attributes {
		storeAttr := RecordForModel(AttributeModel)
		storeAttr.Set("name", aKey)
		storeAttr.Set("id", attr.ID)
		storeAttr.SetFK("model", m.ID)
		storeAttr.SetFK("datatype", attr.Datatype.ID)
		tx.h = tx.h.Insert(storeAttr)
	}

	for _, rel := range m.RightRelationships {
		saveRel(tx, rel)
	}
	for _, rel := range m.LeftRelationships {
		saveRel(tx, rel)
	}

	// done for side effect of gob registration
	RecordForModel(m)
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
