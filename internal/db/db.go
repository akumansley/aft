package db

import (
	"sync"
)

func New() DB {
	appDB := holdDB{
		h:         NewHold(),
		attrs:     map[ID]AttributeLoader{},
		rels:      map[ID]RelationshipLoader{},
		runtimes:  map[ID]FunctionLoader{},
		datatypes: map[ID]DatatypeLoader{}}
	appDB.AddMetaModel()
	return &appDB
}

func NewTest() DB {
	return New()
}

func (db *holdDB) AddMetaModel() {
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

	enums := []Literal{
		StoredAs,
		FunctionSignature,
	}

	for _, e := range enums {
		db.AddLiteral(e)
	}

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

	db.RegisterAttributeLoader(ConcreteAttributeLoader{})
	db.RegisterRelationshipLoader(ConcreteRelationshipLoader{})
	db.RegisterRelationshipLoader(ReverseRelationshipLoader{})
	db.RegisterDatatypeLoader(CoreDatatypeLoader{})
	db.RegisterDatatypeLoader(EnumDatatypeLoader{})

	models := []Literal{
		ModelModel,
		EnumValueModel,
	}

	for _, m := range models {
		db.AddLiteral(m)
	}

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
	RegisterDatatypeLoader(DatatypeLoader)
}

type holdDB struct {
	sync.RWMutex
	h         *Hold
	runtimes  map[ID]FunctionLoader
	attrs     map[ID]AttributeLoader
	rels      map[ID]RelationshipLoader
	datatypes map[ID]DatatypeLoader
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
	db.runtimes[m.ID()] = r
}

func (db *holdDB) RegisterAttributeLoader(l AttributeLoader) {
	m := l.ProvideModel()
	db.AddLiteral(m)
	db.attrs[m.ID()] = l
}

func (db *holdDB) RegisterRelationshipLoader(l RelationshipLoader) {
	m := l.ProvideModel()
	db.AddLiteral(m)
	db.rels[m.ID()] = l
}

func (db *holdDB) RegisterDatatypeLoader(l DatatypeLoader) {
	m := l.ProvideModel()
	db.AddLiteral(m)
	db.datatypes[m.ID()] = l
}

func (db *holdDB) AddLiteral(lit Literal) {
	tx := db.NewRWTx()
	recs, links := lit.MarshalDB()
	for _, rec := range recs {
		tx.Insert(rec)
	}
	for _, link := range links {
		tx.Connect(link.from, link.to, link.rel.ID())
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
		lok := leftI.Next()
		rok := rightI.Next()
		if lok != rok {
			return false
		}
		if lok {
			lv := leftI.Value()
			rv := rightI.Value()
			litem := lv.(item)
			ritem := rv.(item)
			if string(litem.k) != string(ritem.k) {
				return false
			}

			lr, ok := litem.v.(Record)
			if ok {
				rr := ritem.v.(Record)
				if !lr.DeepEquals(rr) {
					return false
				}
			}

		} else {
			return true
		}
	}
}
