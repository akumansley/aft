package db

import (
	"context"
	"sync"

	"awans.org/aft/internal/bus"
)

func New(eb *bus.EventBus) DB {
	if eb == nil {
		eb = bus.New()
	}
	builder := NewBuilder()
	appDB := holdDB{
		b:         builder,
		bus:       eb,
		h:         NewHold(),
		attrs:     map[ID]AttributeLoader{},
		rels:      map[ID]RelationshipLoader{},
		runtimes:  map[ID]FunctionLoader{},
		datatypes: map[ID]DatatypeLoader{},
		ifaces:    map[ID]InterfaceLoader{}}
	mh := MakeAutomigrateHandler(builder)
	eb.RegisterHandler(mh)
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
		typeValidator,
	}

	for _, f := range funcs {
		nr.Save(db.b, f)
	}

	enums := []Literal{
		StoredAs,
	}

	for _, e := range enums {
		db.injectLiteral(e)
	}

	core := []Literal{
		Bool,
		Int,
		String,
		UUID,
		Float,
		Bytes,
		Type,
	}

	for _, d := range core {
		db.injectLiteral(d)
	}

	rwtx := db.NewRWTx()
	db.injectLiteral(GlobalIDAttribute)
	db.injectLiteral(GlobalTypeAttribute)
	rwtx.Connect(GlobalIDAttribute.ID(), UUID.ID(), ConcreteAttributeDatatype.ID())
	rwtx.Connect(GlobalTypeAttribute.ID(), Type.ID(), ConcreteAttributeDatatype.ID())
	rwtx.Commit()

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
		db.injectLiteral(m)
	}
	rwtx = db.makeTx(context.Background())
	rwtx.addImplements(ModelModel.ID(), InterfaceInterface.ID())
	rwtx.addImplements(InterfaceModel.ID(), InterfaceInterface.ID())

	rwtx.addImplements(ConcreteRelationshipModel.ID(), RelationshipInterface.ID())
	rwtx.addImplements(ReverseRelationshipModel.ID(), RelationshipInterface.ID())

	rwtx.addImplements(CoreDatatypeModel.ID(), DatatypeInterface.ID())
	rwtx.addImplements(EnumModel.ID(), DatatypeInterface.ID())

	rwtx.addImplements(NativeFunctionModel.ID(), FunctionInterface.ID())
	rwtx.Commit()
}

type DB interface {
	NewTx() Tx
	NewRWTx() RWTx
	NewTxWithContext(context.Context) Tx
	NewRWTxWithContext(context.Context) RWTx
	DeepEquals(DB) bool
	Iterator() KVIterator
	Builder() *Builder

	AddLiteral(RWTx, Literal)
	RegisterRuntime(FunctionLoader)
	RegisterAttributeLoader(AttributeLoader)
	RegisterRelationshipLoader(RelationshipLoader)
	RegisterDatatypeLoader(DatatypeLoader)
	RegisterNativeFunction(NativeFunctionL)
}

type holdDB struct {
	sync.RWMutex
	b         *Builder
	h         *hold
	runtimes  map[ID]FunctionLoader
	attrs     map[ID]AttributeLoader
	rels      map[ID]RelationshipLoader
	datatypes map[ID]DatatypeLoader
	ifaces    map[ID]InterfaceLoader
	bus       *bus.EventBus
}

func (db *holdDB) Builder() *Builder {
	return db.b
}

func (db *holdDB) NewTxWithContext(ctx context.Context) Tx {
	db.RLock()
	tx := holdTx{initH: db.h, h: db.h, db: db, rw: false, ctx: ctx}
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
	tx := holdTx{initH: db.h, h: db.h, db: db, rw: true, ctx: ctx}
	db.RUnlock()
	return &tx
}

func (db *holdDB) RegisterInterfaceLoader(l InterfaceLoader) {
	m := l.ProvideModel()
	db.injectLiteral(m)
	db.ifaces[m.ID()] = l
}

func (db *holdDB) RegisterRuntime(r FunctionLoader) {
	m := r.ProvideModel()
	db.injectLiteral(m)
	db.runtimes[m.ID()] = r
}

func (db *holdDB) RegisterAttributeLoader(l AttributeLoader) {
	m := l.ProvideModel()
	db.injectLiteral(m)
	db.attrs[m.ID()] = l
}

func (db *holdDB) RegisterRelationshipLoader(l RelationshipLoader) {
	m := l.ProvideModel()
	db.injectLiteral(m)
	db.rels[m.ID()] = l
}

func (db *holdDB) RegisterDatatypeLoader(l DatatypeLoader) {
	m := l.ProvideModel()
	db.injectLiteral(m)
	db.datatypes[m.ID()] = l
}

func (db *holdDB) AddLiteral(tx RWTx, lit Literal) {
	recs, links := lit.MarshalDB(db.b)

	for _, rec := range recs {
		err := tx.Insert(rec)
		if err != nil {
			panic(err)
		}
	}
	for _, link := range links {
		err := tx.Connect(link.From, link.To, link.Rel.ID())
		if err != nil {
			panic(err)
		}
	}
}

func (db *holdDB) injectLiteral(lit Literal) {
	tx := db.makeTx(context.Background())
	recs, links := lit.MarshalDB(db.b)

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
		nr.Save(db.b, nf)
	}

}

func (db *holdDB) Iterator() KVIterator {
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
			lk := leftI.Key()
			rk := rightI.Key()
			if string(lk) != string(rk) {
				return false
			}

			lr, ok := lv.(Record)
			if ok {
				rr := rv.(Record)
				if !lr.DeepEquals(rr) {
					return false
				}
			}

		} else {
			return true
		}
	}
}
