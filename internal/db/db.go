package db

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sort"
	"strings"
	"sync"
)

var (
	ErrData         = errors.New("data-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

func New(ex CodeExecutor) DB {
	appDB := holdDB{h: NewHold(), ex: ex}
	appDB.AddMetaModel()
	return &appDB
}

//tests only rely on golang execution
func NewTest() DB {
	return New(&bootstrapCodeExecutor{})
}

func (db *holdDB) AddMetaModel() {
	var err error
	tx := db.NewRWTx()
	//Add datatypes, enum values and native code
	for _, v := range enumMap {
		r := RecordForModel(EnumValueModel)
		err := SaveEnum(r, v)
		if err != nil {
			panic(err)
		}
		tx.Insert(r)
	}
	//Add native datatypes and their code execution to the tree. Comes before models.
	for _, v := range codeMap {
		err = SaveCode(tx, v)
		if err != nil {
			panic(err)
		}
	}
	for _, v := range datatypeMap {
		err = SaveDatatype(tx, v)
		if err != nil {
			panic(err)
		}
	}
	err = tx.SaveModel(ModelModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(AttributeModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(ModelAttributes)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(RelationshipModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(RelationshipSource)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(RelationshipTarget)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(DatatypeModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(CodeModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(EnumValueModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(DatatypeEnumValues)
	if err != nil {
		panic(err)
	}

	err = tx.SaveRelationship(AttributeDatatype)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(DatatypeValidator)
	if err != nil {
		panic(err)
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
}

type Tx interface {
	Schema() *Schema

	GetRelatedOne(ID, Relationship) (Record, error)
	GetRelatedMany(ID, Relationship) ([]Record, error)
	FindOne(ModelID, Matcher) (Record, error)
	FindMany(ModelID, Matcher) ([]Record, error)
	Ref(ModelID) ModelRef
	Query(ModelRef) Q
}

type RWTx interface {
	// remove
	Schema() *Schema

	// reads
	GetRelatedOne(ID, Relationship) (Record, error)
	GetRelatedMany(ID, Relationship) ([]Record, error)
	GetRelatedManyReverse(ID, Relationship) ([]Record, error)
	FindOne(ModelID, Matcher) (Record, error)
	FindMany(ModelID, Matcher) ([]Record, error)
	Ref(ModelID) ModelRef
	Query(ModelRef) Q

	// writes
	MakeRecord(ModelID) (Record, error)
	Insert(Record) error
	Update(oldRec, newRec Record) error
	Delete(Record) error
	Connect(source, target ID, rel Relationship) error

	Commit() error
}

type holdDB struct {
	sync.RWMutex
	h  *Hold
	ex CodeExecutor
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

func (tx *holdTx) FindOne(modelID ModelID, matcher Matcher) (rec Record, err error) {
	rec, err = tx.h.FindOne(modelID, matcher)
	return
}

func (tx *holdTx) FindMany(modelID ModelID, matcher Matcher) (recs []Record, err error) {
	recs, err = tx.h.FindMany(modelID, matcher)
	return
}

func (tx *holdTx) getRelatedOne(rec Record, rel Relationship) {

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

func (tx *holdTx) Connect(source, target ID, rel Relationship) error {
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

func LoadRel(tx *holdTx, storeRel Record) (rel Relationship, err error) {
	relIf, ok := tx.cache[storeRel.ID()]
	if ok {
		return relIf.(Relationship), err
	}

	if storeRel == nil {
		panic("hi")
	}
	ew := NewRecordWriter(storeRel)

	storeSource, err := tx.h.GetLinkedOne(storeRel, RelationshipSource)
	if err != nil {
		return
	}
	source, err := loadModel(tx, storeSource)
	if err != nil {
		return
	}
	storeTarget, err := tx.h.GetLinkedOne(storeRel, RelationshipTarget)
	if err != nil {
		return
	}
	target, err := loadModel(tx, storeTarget)
	if err != nil {
		return
	}

	// don't initialize source/target yet
	r := Relationship{
		ID:     storeRel.ID(),
		Source: source,
		Target: target,
		Name:   ew.Get("name").(string),
		Multi:  ew.Get("multi").(bool),
	}
	if ew.err != nil {
		return Relationship{}, ew.err
	}
	tx.cache[r.ID] = r
	return r, nil
}

func loadModel(tx *holdTx, storeModel Record) (m Model, err error) {
	mIf, ok := tx.cache[storeModel.ID()]
	if ok {
		return mIf.(Model), err
	}

	ew := NewRecordWriter(storeModel)
	m = Model{
		ID:   ModelID(storeModel.ID()),
		Name: ew.Get("name").(string),
	}
	if ew.err != nil {
		return m, ew.err
	}
	attrs := []Attribute{}

	ami, err := tx.h.GetLinkedMany(storeModel, ModelAttributes)
	if err != nil {
		return
	}

	for _, storeAttr := range ami {
		var storeDatatype Record
		storeDatatype, err = tx.h.GetLinkedOne(storeAttr, AttributeDatatype)
		if err != nil {
			return
		}
		var enum interface{}
		enum, err = storeDatatype.Get("enum")
		if err != nil {
			return Model{}, err
		}
		var native interface{}
		native, err = storeDatatype.Get("native")
		if err != nil {
			return Model{}, err
		}
		var d Datatype
		if enum == true {
			var e Enum
			d, err = e.RecordToStruct(storeDatatype, tx)
			if err != nil {
				return Model{}, err
			}
		} else if native == true {
			var c coreDatatype
			d, err = c.RecordToStruct(storeDatatype, tx)
			if err != nil {
				return Model{}, err
			}
		} else {
			var c DatatypeStorage
			d, err = c.RecordToStruct(storeDatatype, tx)
			if err != nil {
				return Model{}, err
			}
		}
		name, err := storeAttr.Get("name")
		if err != nil {
			return Model{}, err
		}
		attr := Attribute{
			Datatype: d,
			ID:       storeAttr.ID(),
			Name:     name.(string),
		}
		attrs = append(attrs, attr)
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Name < attrs[j].Name
	})
	m.Attributes = attrs

	tx.cache[ID(m.ID)] = m
	return m, nil
}

func (tx *holdTx) GetModelByID(id ModelID) (m Model, err error) {

	storeModel, err := tx.h.FindOne(ModelModel.ID, EqID(ID(id)))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, id)
	}
	return loadModel(tx, storeModel)
}

func (tx *holdTx) GetModel(modelName string) (m Model, err error) {
	modelName = strings.ToLower(modelName)
	storeModel, err := tx.h.FindOne(ModelModel.ID, Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	return loadModel(tx, storeModel)
}

func (tx *holdTx) GetRelationship(id ID) (rel Relationship, err error) {
	storeRel, err := tx.FindOne(RelationshipModel.ID, EqID(id))
	if err != nil {
		return
	}
	rel, err = LoadRel(tx, storeRel)
	return
}

func (tx *holdTx) GetRelationships(model Model) (rels []Relationship, err error) {
	recs, err := tx.h.GetLinkedManyReverseByID(ID(model.ID), RelationshipSource)
	if err != nil {
		return
	}
	for _, r := range recs {
		var rel Relationship
		rel, err = LoadRel(tx, r)
		if err != nil {
			return
		}
		rels = append(rels, rel)
	}
	return
}

// TODO rewrite
func SaveDatatype(tx RWTx, d Datatype) error {
	storeDatatype := RecordForModel(DatatypeModel)
	d.FillRecord(storeDatatype)

	tx.Insert(storeDatatype)
	// tx.Connect(storeDatatype, storeValidator, DatatypeValidator)

	return nil
}

func SaveCode(tx RWTx, c Code) error {
	storeCode := RecordForModel(CodeModel)
	ew := NewRecordWriter(storeCode)
	ew.Set("id", uuid.UUID(c.ID))
	ew.Set("name", c.Name)
	ew.Set("runtime", uuid.UUID(c.Runtime.ID))
	ew.Set("functionSignature", uuid.UUID(c.FunctionSignature.ID))
	ew.Set("code", c.Code)
	return ew.err
}

func SaveEnum(storeEnum Record, e EnumValue) error {
	ew := NewRecordWriter(storeEnum)
	ew.Set("id", uuid.UUID(e.ID))
	ew.Set("name", e.Name)
	// ew.SetFK("datatype", e.Datatype)
	// if ew.err != nil {
	// 	return ew.err
	// }
	// tx.Insert(storeCode)
	return nil
}

func (tx *holdTx) SaveRelationship(rel Relationship) (err error) {
	storeRel := RecordForModel(RelationshipModel)
	ew := NewRecordWriter(storeRel)
	ew.Set("name", rel.Name)
	ew.Set("multi", rel.Multi)
	ew.Set("id", uuid.UUID(rel.ID))
	if ew.err != nil {
		return ew.err
	}
	tx.h = tx.h.Insert(storeRel)

	storeSource, err := tx.FindOne(ModelModel.ID, EqID(ID(rel.Source.ID)))
	if err != nil {
		return
	}
	tx.Connect(storeRel, storeSource, RelationshipSource)

	storeTarget, err := tx.FindOne(ModelModel.ID, EqID(ID(rel.Target.ID)))
	if err != nil {
		return
	}
	tx.Connect(storeRel, storeTarget, RelationshipTarget)
	return
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m Model) error {
	tx.ensureWrite()
	storeModel := RecordForModel(ModelModel)
	ew := NewRecordWriter(storeModel)
	ew.Set("name", m.Name)
	ew.Set("id", uuid.UUID(m.ID))
	tx.h = tx.h.Insert(storeModel)
	if ew.err != nil {
		return ew.err
	}
	for _, attr := range m.Attributes {
		storeAttr := RecordForModel(AttributeModel)
		ew = NewRecordWriter(storeAttr)
		ew.Set("name", attr.Name)
		ew.Set("id", uuid.UUID(attr.ID))
		ew.Set("datatypeID", uuid.UUID(attr.Datatype.GetID())) //TODO remove hack
		storeDatatype, _ := tx.FindOne(DatatypeModel.ID, EqID(attr.Datatype.GetID()))
		tx.Connect(storeModel, storeAttr, ModelAttributes)
		tx.Connect(storeAttr, storeDatatype, AttributeDatatype)
		tx.h = tx.h.Insert(storeAttr)
	}

	// done for side effect of gob registration
	RecordForModel(m)
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func (tx *holdTx) MakeRecord(modelID ModelID) (rec Record, err error) {
	m, err := tx.GetModelByID(modelID)
	rec = RecordForModel(m)
	return
}

func (tx *holdTx) Commit() error {
	tx.ensureWrite()
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()
	return nil
}

type errWriter struct {
	r   Record
	err error
}

func NewRecordWriter(r Record) *errWriter {
	return &errWriter{r, nil}
}

func (ew *errWriter) Set(key string, val interface{}) {
	if ew.err == nil {
		ew.err = ew.r.Set(key, val)
	}
}

func (ew *errWriter) Get(key string) interface{} {
	var out interface{}
	if ew.err == nil {
		out, ew.err = ew.r.Get(key)
		return out
	}
	return nil
}
