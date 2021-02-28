package db

import (
	"context"
	"fmt"
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
	db.injectLiteral(nr.ProvideModel())

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
		db.injectLiteral(f)
	}

	core := []DatatypeL{
		StoredAs,
		FuncType,
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

	// bootstrap id/type
	db.injectLiteral(GlobalIDAttribute)
	db.injectLiteral(GlobalTypeAttribute)

	rwtx := db.NewRWTx()
	rwtx.Connect(GlobalIDAttribute.ID(), UUID.ID(), ConcreteAttributeDatatype.ID())
	rwtx.Connect(GlobalTypeAttribute.ID(), Type.ID(), ConcreteAttributeDatatype.ID())
	rwtx.Commit()

	db.RegisterInterfaceLoader(ModelInterfaceLoader{})
	db.RegisterInterfaceLoader(InterfaceInterfaceLoader{})

	db.RegisterAttributeLoader(ConcreteAttributeLoader{})

	db.RegisterRelationshipLoader(InterfaceRelationshipLoader{})
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
		ModuleModel,
		ModuleInterfaces,
		ModuleFunctions,
		ModuleDatatypes,
	}

	for _, m := range models {
		db.injectLiteral(m)
	}
	rwtx = db.NewRWTx()
	rwtx.addImplements(ModelModel.ID(), InterfaceInterface.ID())
	rwtx.addImplements(InterfaceModel.ID(), InterfaceInterface.ID())

	rwtx.addImplements(InterfaceRelationshipModel.ID(), RelationshipInterface.ID())
	rwtx.addImplements(ConcreteRelationshipModel.ID(), RelationshipInterface.ID())
	rwtx.addImplements(ReverseRelationshipModel.ID(), RelationshipInterface.ID())

	rwtx.addImplements(CoreDatatypeModel.ID(), DatatypeInterface.ID())
	rwtx.addImplements(EnumModel.ID(), DatatypeInterface.ID())

	rwtx.addImplements(NativeFunctionModel.ID(), FunctionInterface.ID())
	rwtx.Commit()

	dbMod := MakeModule(
		MakeID("fde9bb79-2b8e-4dd4-b830-140368606d57"),
		"db",
		"awans.org/aft/internal/db",
		[]InterfaceL{
			ModelModel,
			InterfaceModel,
			InterfaceRelationshipModel,
			ConcreteRelationshipModel,
			ReverseRelationshipModel,
			CoreDatatypeModel,
			EnumModel,
			NativeFunctionModel,
			EnumValueModel,
			InterfaceInterface,
			RelationshipInterface,
			FunctionInterface,
			DatatypeInterface,
			ModuleModel,
		},
		[]FunctionL{
			boolValidator,
			intValidator,
			stringValidator,
			uuidValidator,
			floatValidator,
			bytesValidator,
			typeValidator,
		},
		core,
		nil,
	)
	db.AddLiteral(dbMod)
}

type DB interface {
	NewTx() Tx
	NewRWTx() RWTx
	NewTxWithContext(context.Context) Tx
	NewRWTxWithContext(context.Context) RWTx
	DeepEquals(DB) bool
	Iterator() KVIterator
	Builder() *Builder

	AddLiteral(Literal)
	RegisterRuntime(FunctionLoader)
	RegisterAttributeLoader(AttributeLoader)
	RegisterRelationshipLoader(RelationshipLoader)
	RegisterDatatypeLoader(DatatypeLoader)
	RegisterNativeFunction(NativeFunctionLiteral)
}

type holdDB struct {
	sync.RWMutex
	writer    sync.Mutex
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
	hSnap := db.h
	db.RUnlock()
	holdTx := &holdTx{initH: hSnap, h: hSnap, db: db, rw: false}

	tx := &txWithContext{holdTx: holdTx}
	ctx = WithTx(ctx, tx)
	tx.ctx = ctx

	return tx
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

func (db *holdDB) makeTx(ctx context.Context) *txWithContext {
	db.writer.Lock()
	db.RLock()
	hSnap := db.h
	db.RUnlock()
	holdTx := &holdTx{initH: hSnap, h: hSnap, db: db, rw: true}

	tx := &txWithContext{holdTx: holdTx}
	ctx = WithRWTx(ctx, tx)
	tx.ctx = ctx

	return tx
}

func (db *holdDB) RegisterInterfaceLoader(l InterfaceLoader) {
	m := l.ProvideModel()
	db.injectLiteral(m)
	db.ifaces[m.ID()] = l
}

func (db *holdDB) RegisterRuntime(r FunctionLoader) {
	m := r.ProvideModel()
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

func (db *holdDB) AddLiteral(lit Literal) {
	tx := db.NewRWTx()
	recs, links := lit.MarshalDB(db.b)

	for _, rec := range recs {
		err := tx.Insert(rec)
		if err != nil {
			err = fmt.Errorf("%w adding literal %v\n", err, lit.ID())
			panic(err)
		}
	}
	for _, link := range links {
		relID := link.Rel.ID()
		ir, ok := link.Rel.(InterfaceRelationshipL)
		if ok {
			concreteSourceID := link.From.InterfaceID()
			rel := tx.Ref(RelationshipInterface.ID())
			source := tx.Ref(InterfaceInterface.ID())
			conRel, err := tx.Query(rel,
				Filter(rel, Eq("name", ir.Name_)),
				Join(source, rel.Rel(ConcreteRelationshipSource.Load(tx))),
				Filter(source, EqID(concreteSourceID)),
			).OneRecord()
			if err != nil {
				panic(err)
			}
			relID = conRel.ID()
		}
		err := tx.Connect(link.From.ID(), link.To.ID(), relID)
		if err != nil {
			err = fmt.Errorf("%w adding literal %v\n", err, lit.ID())
			panic(err)
		}
	}
	tx.Commit()
}

func (db *holdDB) injectLiteral(lit Literal) {
	tx := db.makeTx(context.Background())
	recs, links := lit.MarshalDB(db.b)

	for _, rec := range recs {
		tx.h = tx.h.Insert(rec)
	}
	for _, link := range links {
		tx.h = tx.h.Link(link.From.ID(), link.To.ID(), link.Rel.ID())
	}

	tx.Commit()
}

func (db *holdDB) RegisterNativeFunction(nf NativeFunctionLiteral) {
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
