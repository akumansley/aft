package db

import (
	"context"
	"sync"

	"awans.org/aft/internal/bus"
)

func New(b *bus.EventBus) DB {
	if b == nil {
		b = bus.New()
	}
	appDB := holdDB{
		bus:       b,
		h:         NewHold(),
		attrs:     map[ID]AttributeLoader{},
		rels:      map[ID]RelationshipLoader{},
		runtimes:  map[ID]FunctionLoader{},
		datatypes: map[ID]DatatypeLoader{},
		ifaces:    map[ID]InterfaceLoader{}}
	appDB.AddMetaModel()
	return &appDB
}

func NewTest() DB {
	return New(nil)
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
		bytesValidator,
	}

	for _, f := range funcs {
		nr.Save(f)
	}

	enums := []Literal{
		StoredAs,
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
		Bytes,
	}

	for _, d := range core {
		db.AddLiteral(d)
	}

	db.RegisterInterfaceLoader(ModelInterfaceLoader{})
	db.RegisterInterfaceLoader(InterfaceInterfaceLoader{})

	db.RegisterAttributeLoader(ConcreteAttributeLoader{})

	db.RegisterRelationshipLoader(ConcreteRelationshipLoader{})
	db.RegisterRelationshipLoader(ReverseRelationshipLoader{})

	db.RegisterDatatypeLoader(CoreDatatypeLoader{})
	db.RegisterDatatypeLoader(EnumDatatypeLoader{})

	models := []Literal{
		EnumValueModel,
		InterfaceInterface,
		RelationshipInterface,
		FunctionInterface,
		DatatypeInterface,
	}

	for _, m := range models {
		db.AddLiteral(m)
	}

}

type DB interface {
	NewTx() Tx
	NewRWTx() RWTx
	NewTxWithContext(context.Context) Tx
	NewRWTxWithContext(context.Context) RWTx
	DeepEquals(DB) bool
	Iterator() Iterator

	AddLiteral(Literal)
	RegisterRuntime(FunctionLoader)
	RegisterAttributeLoader(AttributeLoader)
	RegisterRelationshipLoader(RelationshipLoader)
	RegisterDatatypeLoader(DatatypeLoader)
	RegisterNativeFunction(NativeFunctionL)
}

type holdDB struct {
	sync.RWMutex
	h         *Hold
	runtimes  map[ID]FunctionLoader
	attrs     map[ID]AttributeLoader
	rels      map[ID]RelationshipLoader
	datatypes map[ID]DatatypeLoader
	ifaces    map[ID]InterfaceLoader
	bus       *bus.EventBus
}

func (db *holdDB) NewTxWithContext(ctx context.Context) Tx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: false}
	db.RUnlock()
	return &tx
}

func (db *holdDB) NewTx() Tx {
	return db.NewTxWithContext(context.Background())
}

func (db *holdDB) NewRWTx() RWTx {
	return db.NewRWTxWithContext(context.Background())
}

func (db *holdDB) NewRWTxWithContext(ctx context.Context) RWTx {
	return db.makeTx(ctx)
}

func (db *holdDB) makeTx(ctx context.Context) *holdTx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: true}
	db.RUnlock()
	return &tx
}

func (db *holdDB) RegisterInterfaceLoader(l InterfaceLoader) {
	m := l.ProvideModel()
	db.AddLiteral(m)
	db.ifaces[m.ID()] = l
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
	tx := db.makeTx(context.Background())
	recs, links := lit.MarshalDB()

	// circumvent proper Tx checks for bootstrapping;
	// danger!
	for _, rec := range recs {
		tx.h = tx.h.Insert(rec)
	}
	for _, link := range links {
		tx.h = tx.h.Link(link.From, link.To, link.Rel.ID())
	}
	tx.Commit()
}

func (db *holdDB) RegisterNativeFunction(nf NativeFunctionL) {
	fl, ok := db.runtimes[NativeFunctionModel.ID()]
	if !ok {
		panic("bad order")
	}

	if nr, ok := fl.(*NativeRuntime); ok {
		nr.Save(nf)
	}

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
