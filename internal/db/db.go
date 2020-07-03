package db

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	for _, v := range datatypeMap {
		r := RecordForModel(DatatypeModel)
		err := SaveDatatype(r, v)
		if err != nil {
			panic(err)
		}
		tx.Insert(r)
	}
	for _, v := range codeMap {
		r := RecordForModel(CodeModel)
		err := SaveCode(r, v)
		if err != nil {
			panic(err)
		}
		tx.Insert(r)
	}
	err := tx.SaveModel(ModelModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(AttributeModel)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(RelationshipModel)
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
	Ex() CodeExecutor
	GetModel(string) (Model, error)
	GetModelByID(ModelID) (Model, error)
	MakeRecord(ModelID) (Record, error)
	FindOne(ModelID, Matcher) (Record, error)
	FindMany(ModelID, Matcher) ([]Record, error)
	Ref(ModelID) ModelRef
	Query(ModelRef) Q
}

type RWTx interface {
	Ex() CodeExecutor
	// remove
	GetModel(string) (Model, error)
	GetModelByID(ModelID) (Model, error)
	SaveModel(Model) error

	FindOne(ModelID, Matcher) (Record, error)
	FindMany(ModelID, Matcher) ([]Record, error)
	MakeRecord(ModelID) (Record, error)
	Ref(ModelID) ModelRef
	Query(ModelRef) Q

	// these are good, i think
	Insert(Record) error
	Connect(from, to Record, fromRel Relationship) error
	Update(oldRec, newRec Record) error
	Delete(Record) error

	Commit() error
}

type holdDB struct {
	sync.RWMutex
	h  *Hold
	ex CodeExecutor
}

type holdTx struct {
	h  *Hold
	db *holdDB
	rw bool
	ex CodeExecutor
}

func (tx *holdTx) ensureWrite() {
	if !tx.rw {
		panic("Tried to write in a read only tx")
	}
}

func (db *holdDB) NewTx() Tx {
	db.RLock()
	tx := holdTx{h: db.h, rw: false, ex: db.ex}
	db.RUnlock()
	return &tx
}

func (db *holdDB) NewRWTx() RWTx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: true, ex: db.ex}
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

func (tx *holdTx) Ex() CodeExecutor {
	return tx.ex
}

func (tx *holdTx) FindOne(modelID ModelID, matcher Matcher) (rec Record, err error) {
	rec, err = tx.h.FindOne(modelID, matcher)
	return
}

func (tx *holdTx) FindMany(modelID ModelID, matcher Matcher) (recs []Record, err error) {
	recs, err = tx.h.FindMany(modelID, matcher)
	return
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

func (tx *holdTx) Connect(left, right Record, rel Relationship) error {
	tx.ensureWrite()

	if rel.LeftBinding == BelongsTo && (rel.RightBinding == HasOne || rel.RightBinding == HasMany) {
		// FK left
		err := left.SetFK(rel.LeftName, right.ID())
		if err != nil {
			return err
		}
	} else if rel.RightBinding == BelongsTo && (rel.LeftBinding == HasOne || rel.LeftBinding == HasMany) {
		// FK right
		err := right.SetFK(rel.RightName, left.ID())
		if err != nil {
			return err
		}
	} else if rel.LeftBinding == HasManyAndBelongsToMany && rel.RightBinding == HasManyAndBelongsToMany {
		// Join table
		panic("Many to many relationships not implemented yet")
	} else {
		return fmt.Errorf("Trying to connect invalid relationship")
	}
	h1 := tx.h.Insert(left)
	h2 := h1.Insert(right)
	tx.h = h2
	return nil
}

func (tx *holdTx) Delete(rec Record) error {
	tx.ensureWrite()
	tx.h = tx.h.Delete(rec)
	return nil
}

func LoadRel(storeRel Record) (Relationship, error) {
	ew := NewRecordWriter(storeRel)
	r := Relationship{
		ID:           storeRel.ID(),
		LeftBinding:  RelType(ew.Get("leftBinding").(int64)),
		LeftModelID:  ModelID(ew.GetFK("leftModel")),
		LeftName:     ew.Get("leftName").(string),
		RightBinding: RelType(ew.Get("rightBinding").(int64)),
		RightModelID: ModelID(ew.GetFK("rightModel")),
		RightName:    ew.Get("rightName").(string),
	}
	if ew.err != nil {
		return Relationship{}, ew.err
	}
	return r, nil
}

func loadModel(tx *holdTx, storeModel Record) (Model, error) {
	ew := NewRecordWriter(storeModel)
	m := Model{
		ID:     ModelID(storeModel.ID()),
		Name:   ew.Get("name").(string),
		System: ew.Get("system").(bool),
	}
	if ew.err != nil {
		return Model{}, nil
	}
	attrs := []Attribute{}
	// make ModelID a dynamic key
	ami, err := tx.FindMany(AttributeModel.ID, EqFK("model", ID(m.ID)))
	if err != nil {
		return Model{}, err
	}
	for _, storeAttr := range ami {
		dk, err := storeAttr.GetFK("datatype")
		if err != nil {
			return Model{}, err
		}
		storeDatatype, err := tx.h.FindOne(DatatypeModel.ID, EqID(dk))
		if err != nil {
			return Model{}, err
		}
		enum, err := storeDatatype.Get("enum")
		if err != nil {
			return Model{}, err
		}
		native, err := storeDatatype.Get("native")
		if err != nil {
			return Model{}, err
		}
		var d Datatype
		if enum == true {
			var e Enum
			d, err = e.RecordToDatatype(storeDatatype, tx)
			if err != nil {
				return Model{}, err
			}
		} else if native == true {
			var c coreDatatype
			d, err = c.RecordToDatatype(storeDatatype, tx)
			if err != nil {
				return Model{}, err
			}
		} else {
			var c DatatypeStorage
			d, err = c.RecordToDatatype(storeDatatype, tx)
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
	m.Attributes = attrs

	lRels := []Relationship{}
	rmi, err := tx.FindMany(RelationshipModel.ID, EqFK("leftModel", ID(m.ID)))
	if err != nil {
		return Model{}, err
	}
	for _, storeRel := range rmi {
		rel, err := LoadRel(storeRel)
		if err != nil {
			return Model{}, err
		}
		lRels = append(lRels, rel)
	}
	m.LeftRelationships = lRels

	rRels := []Relationship{}
	rmi, err = tx.FindMany(RelationshipModel.ID, EqFK("rightModel", ID(m.ID)))
	if err != nil {
		return Model{}, err
	}
	for _, storeRel := range rmi {
		rel, err := LoadRel(storeRel)
		if err != nil {
			return Model{}, err
		}
		rRels = append(rRels, rel)
	}
	m.RightRelationships = rRels
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
	storeModel, err := tx.h.FindOne(ModelModel.ID, Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	return loadModel(tx, storeModel)
}

func SaveDatatype(storeDatatype Record, d Datatype) error {
	return d.FillRecord(storeDatatype)
}

func SaveEnum(storeEnum Record, e EnumValue) error {
	ew := NewRecordWriter(storeEnum)
	ew.Set("id", uuid.UUID(e.ID))
	ew.Set("name", e.Name)
	ew.SetFK("datatype", e.Datatype)
	return ew.err
}

func SaveRel(rel Relationship) (Record, error) {
	storeRel := RecordForModel(RelationshipModel)
	ew := NewRecordWriter(storeRel)
	ew.SetFK("leftModel", ID(rel.LeftModelID))
	ew.Set("leftName", rel.LeftName)
	ew.Set("leftBinding", int64(rel.LeftBinding))
	ew.SetFK("rightModel", ID(rel.RightModelID))
	ew.Set("rightName", rel.RightName)
	ew.Set("rightBinding", int64(rel.RightBinding))
	ew.Set("id", uuid.UUID(rel.ID))
	if ew.err != nil {
		return storeRel, ew.err
	}
	return storeRel, nil
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m Model) error {
	tx.ensureWrite()
	storeModel := RecordForModel(ModelModel)
	ew := NewRecordWriter(storeModel)
	ew.Set("name", m.Name)
	ew.Set("id", uuid.UUID(m.ID))
	ew.Set("system", m.System)
	tx.h = tx.h.Insert(storeModel)
	if ew.err != nil {
		return ew.err
	}
	for _, attr := range m.Attributes {
		storeAttr := RecordForModel(AttributeModel)
		ew = NewRecordWriter(storeAttr)
		ew.Set("name", attr.Name)
		ew.Set("id", uuid.UUID(attr.ID))
		ew.Set("datatypeId", uuid.UUID(attr.Datatype.GetID())) //TODO remove hack
		ew.SetFK("model", ID(m.ID))
		ew.SetFK("datatype", attr.Datatype.GetID())
		tx.h = tx.h.Insert(storeAttr)
	}

	for _, rel := range m.RightRelationships {
		storeRel, err := SaveRel(rel)
		if err != nil {
			return err
		}
		tx.h = tx.h.Insert(storeRel)
	}
	for _, rel := range m.LeftRelationships {
		storeRel, err := SaveRel(rel)
		if err != nil {
			return err
		}
		tx.h = tx.h.Insert(storeRel)
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

func (ew *errWriter) SetFK(key string, val ID) {
	if ew.err == nil {
		ew.err = ew.r.SetFK(key, val)
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

func (ew *errWriter) GetFK(key string) ID {
	var out ID
	if ew.err == nil {
		out, ew.err = ew.r.GetFK(key)
		return out
	}
	return ID(uuid.Nil)
}
