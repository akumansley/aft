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
	MakeRecord(uuid.UUID) Record
	FindOne(uuid.UUID, Matcher) (Record, error)
	FindMany(uuid.UUID, Matcher) ([]Record, error)
	Ref(uuid.UUID) ModelRef
	Query(ModelRef) Q
}

type RWTx interface {
	// remove
	GetModel(string) (Model, error)
	GetModelByID(uuid.UUID) (Model, error)
	SaveModel(Model) error

	FindOne(uuid.UUID, Matcher) (Record, error)
	FindMany(uuid.UUID, Matcher) ([]Record, error)
	MakeRecord(uuid.UUID) Record
	Ref(uuid.UUID) ModelRef
	Query(ModelRef) Q

	// these are good, i think
	Insert(Record) error
	Update(oldRec, newRec Record) error
	Connect(from, to Record, fromRel Relationship) error

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

func (tx *holdTx) FindOne(modelID uuid.UUID, matcher Matcher) (rec Record, err error) {
	rec, err = tx.h.FindOne(modelID, matcher)
	return
}

func (tx *holdTx) FindMany(modelID uuid.UUID, matcher Matcher) (recs []Record, err error) {
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

func LoadRel(storeRel Record) (Relationship, error) {
	ew := &errWriter{}
	ew.NewRecordWriter(storeRel)
	r := Relationship{
		ID:           storeRel.ID(),
		LeftBinding:  RelType(ew.Get("leftBinding").(int64)),
		LeftModelID:  ew.GetFK("leftModel"),
		LeftName:     ew.Get("leftName").(string),
		RightBinding: RelType(ew.Get("rightBinding").(int64)),
		RightModelID: ew.GetFK("rightModel"),
		RightName:    ew.Get("rightName").(string),
	}
	if ew.err != nil {
		return Relationship{}, ew.err
	}
	return r, nil
}

func loadModel(tx *holdTx, storeModel Record) (Model, error) {
	ew := &errWriter{}
	ew.NewRecordWriter(storeModel)
	m := Model{
		ID:   storeModel.ID(),
		Name: ew.Get("name").(string),
	}
	if ew.err != nil {
		return Model{}, nil
	}
	attrs := make(map[string]Attribute)

	// make ModelID a dynamic key
	ami, err := tx.FindMany(AttributeModel.ID, EqFK("model", m.ID))
	if err != nil {
		return Model{}, err
	}
	for _, storeAttr := range ami {
		dk, err := storeAttr.GetFK("datatype")
		if err != nil {
			return Model{}, err
		}
		storeDatatype, _ := tx.h.FindOne(DatatypeModel.ID, Eq("id", dk))
		vk, err := storeDatatype.GetFK("validator")
		if err != nil {
			return Model{}, err
		}
		storeValidator, _ := tx.h.FindOne(CodeModel.ID, Eq("id", vk))
		if ew.err != nil {
			return Model{}, err
		}
		ew.NewRecordWriter(storeValidator)
		validator := Code{
			ID:                storeValidator.ID(),
			Name:              ew.Get("name").(string),
			Runtime:           Runtime(ew.Get("runtime").(int64)),
			Code:              ew.Get("code").(string),
			FunctionSignature: FunctionSignature(ew.Get("functionSignature").(int64)),
			executor:          tx.ex,
		}
		if _, ok := codeMap[validator.ID]; ok {
			validator.Function = codeMap[validator.ID].Function
		}
		if ew.err != nil {
			return Model{}, err
		}
		ew.NewRecordWriter(storeDatatype)
		d := Datatype{
			ID:        storeDatatype.ID(),
			Name:      ew.Get("name").(string),
			Validator: validator,
			StoredAs:  Storage(ew.Get("storedAs").(int64)),
		}
		attr := Attribute{
			Datatype: d,
			ID:       storeAttr.ID(),
		}
		name, err := storeAttr.Get("name")
		if err != nil {
			return Model{}, err
		}
		attrs[name.(string)] = attr
	}
	m.Attributes = attrs

	lRels := []Relationship{}
	rmi, err := tx.FindMany(RelationshipModel.ID, EqFK("leftModel", m.ID))
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
	rmi, err = tx.FindMany(RelationshipModel.ID, EqFK("rightModel", m.ID))
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
	if ew.err != nil {
		return Model{}, err
	}
	return m, nil
}

func (tx *holdTx) GetModelByID(id uuid.UUID) (m Model, err error) {
	storeModel, err := tx.h.FindOne(ModelModel.ID, Eq("id", id))
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

func SaveDatatype(storeDatatype Record, d Datatype) error {
	ew := errWriter{}
	ew.NewRecordWriter(storeDatatype)
	ew.Set("id", d.ID)
	ew.Set("name", d.Name)
	ew.Set("storedAs", int64(d.StoredAs))
	ew.SetFK("validator", d.Validator.ID)
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func SaveCode(storeCode Record, c Code) error {
	ew := errWriter{}
	ew.NewRecordWriter(storeCode)
	ew.Set("id", c.ID)
	ew.Set("name", c.Name)
	ew.Set("runtime", int64(c.Runtime))
	ew.Set("functionSignature", int64(c.FunctionSignature))
	ew.Set("code", c.Code)
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func SaveRel(rel Relationship) (Record, error) {
	ew := errWriter{}
	storeRel := RecordForModel(RelationshipModel)
	ew.NewRecordWriter(storeRel)
	ew.SetFK("leftModel", rel.LeftModelID)
	ew.Set("leftName", rel.LeftName)
	ew.Set("leftBinding", int64(rel.LeftBinding))
	ew.SetFK("rightModel", rel.RightModelID)
	ew.Set("rightName", rel.RightName)
	ew.Set("rightBinding", int64(rel.RightBinding))
	ew.Set("id", rel.ID)
	if ew.err != nil {
		return storeRel, ew.err
	}
	return storeRel, nil
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m Model) error {
	ew := errWriter{}
	tx.ensureWrite()
	storeModel := RecordForModel(ModelModel)
	ew.NewRecordWriter(storeModel)
	ew.Set("name", m.Name)
	ew.Set("id", m.ID)
	tx.h = tx.h.Insert(storeModel)
	if ew.err != nil {
		return ew.err
	}
	for aKey, attr := range m.Attributes {
		storeAttr := RecordForModel(AttributeModel)
		ew.NewRecordWriter(storeAttr)
		ew.Set("name", aKey)
		ew.Set("id", attr.ID)
		ew.Set("datatypeId", attr.Datatype.ID) //TODO remove hack
		ew.SetFK("model", m.ID)
		ew.SetFK("datatype", attr.Datatype.ID)
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

func (tx *holdTx) MakeRecord(modelID uuid.UUID) Record {
	m, _ := tx.GetModelByID(modelID)
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

type errWriter struct {
	r   Record
	err error
}

func (ew *errWriter) NewRecordWriter(r Record) {
	ew.r = r
}
func (ew *errWriter) Set(key string, val interface{}) {
	if ew.err == nil {
		ew.err = ew.r.Set(key, val)
	}
}

func (ew *errWriter) SetFK(key string, val uuid.UUID) {
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

func (ew *errWriter) GetFK(key string) uuid.UUID {
	var out uuid.UUID
	if ew.err == nil {
		out, ew.err = ew.r.GetFK(key)
		return out
	}
	return uuid.Nil
}
