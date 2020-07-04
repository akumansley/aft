package db

import (
	"sync"
)

func New() DB {
	appDB := holdDB{
		h:     NewHold(),
		attrs: map[ID]AttributeLoader{},
		rels:  map[ID]RelationshipLoader{}}
	appDB.AddMetaModel()
	return &appDB
}

//tests only rely on golang execution
func NewTest() DB {
	return New()
}

func (db *holdDB) AddMetaModel() {
	// first add the native runtime
	nr := NewNativeRuntime(db)
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
		db.AddLiteral(d)
	}

	models := []Literal{
		ModelModel,
		ConcreteRelationshipModel,
		CoreDatatypeModel,
		EnumValueModel,
	}

	relationships := []Literal{
		ModelAttributes,
		ConcreteRelationshipSource,
		ConcreteRelationshipTarget,
		ConcreteAttributeDatatype,
		DatatypeValidator,
	}

	for _, m := range models {
		db.AddLiteral(m)
	}
	for _, r := range relationships {
		db.AddLiteral(r)
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
	AddLiteral(Literal)
	RegisterRuntime(FunctionLoader)
	RegisterAttributeLoader(AttributeLoader)
	RegisterRelationshipLoader(RelationshipLoader)
}

type holdDB struct {
	sync.RWMutex
	h        *Hold
	runtimes map[ID]FunctionLoader
	attrs    map[ID]AttributeLoader
	rels     map[ID]RelationshipLoader
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

func (db *holdDB) RegisterRuntime(r FunctionLoader) {
	m := r.ProvideModel()
	db.AddLiteral(m)
	db.runtimes[m.ID] = r
}

func (db *holdDB) RegisterAttributeLoader(l AttributeLoader) {
	m := l.ProvideModel()
	db.AddLiteral(m)
	db.attrs[m.ID] = l
}

func (db *holdDB) RegisterRelationshipLoader(l RelationshipLoader) {
	m := l.ProvideModel()
	db.AddLiteral(m)
	db.rels[m.ID] = l
}

func (db *holdDB) AddLiteral(lit Literal) {
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
