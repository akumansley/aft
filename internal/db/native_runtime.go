package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Func func(context.Context, []interface{}) (interface{}, error)

// Model

var NativeFunctionModel = MakeModel(
	MakeID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	"nativeFunction",
	[]AttributeL{
		nfName,
		nfArity,
		nfType,
	},
	[]RelationshipL{NativeFunctionModule},
	[]ConcreteInterfaceL{FunctionInterface},
)

var nfName = MakeConcreteAttribute(
	MakeID("c47bcd30-01ea-467f-ad02-114342070241"),
	"name",
	String,
)

var nfArity = MakeConcreteAttribute(
	MakeID("dd154c21-2822-41e2-80e3-8489babc907e"),
	"arity",
	Int,
)

var nfType = MakeConcreteAttribute(
	MakeID("a7526d8c-7354-4898-b2c1-2196da0440b7"),
	"funcType",
	FuncType,
)

var NativeFunctionModule = MakeConcreteRelationship(
	MakeID("73fbe1e4-4f60-4031-aedd-ae1f14a4d1e6"),
	"module",
	false,
	ModuleModel,
)

// Loader

type NativeRuntime struct {
	fMap map[ID]Func
	db   DB
}

func NewNativeRuntime(db DB) *NativeRuntime {
	return &NativeRuntime{
		fMap: make(map[ID]Func),
		db:   db,
	}
}

func (nr *NativeRuntime) ProvideModel() ModelL {
	return NativeFunctionModel
}

func (nr *NativeRuntime) Load(rec Record) Function {
	return nativeFunction{rec, nr}
}

func MakeNativeFunction(id ID, name string, arity int, funcType EnumValueL, function Func) NativeFunctionL {
	return NativeFunctionL{
		id, name, arity, funcType, function,
	}
}

// Literal

type NativeFunctionL struct {
	ID_      ID         `record:"id"`
	Name_    string     `record:"name"`
	Arity_   int        `record:"arity"`
	FuncType EnumValueL `record:"funcType"`
	Function Func
}

func (lit NativeFunctionL) ID() ID {
	return lit.ID_
}

func (lit NativeFunctionL) InterfaceID() ID {
	return NativeFunctionModel.ID()
}

func (lit NativeFunctionL) InterfaceName() string {
	return NativeFunctionModel.Name_
}

func (lit NativeFunctionL) Func() Func {
	return lit.Function
}

func (lit NativeFunctionL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	return
}

func (lit NativeFunctionL) Load(tx Tx) Function {
	f, err := tx.Schema().GetFunctionByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return f
}

func (nr *NativeRuntime) Save(b *Builder, lit NativeFunctionLiteral) {
	f := lit.Func()
	tx := nr.db.NewRWTx()
	rec := MarshalRecord(b, lit)
	tx.Insert(rec)
	nr.fMap[lit.ID()] = f
	tx.Commit()
}

// Dynamic

type nativeFunction struct {
	rec Record
	nr  *NativeRuntime
}

func (nf nativeFunction) ID() ID {
	return nf.rec.ID()
}

func (nf nativeFunction) Name() string {
	return nf.rec.MustGet("name").(string)
}

func (nf nativeFunction) Arity() int {
	a := int(nf.rec.MustGet("arity").(int64))
	return a
}

func (nf nativeFunction) FuncType(tx Tx) EnumValue {
	evID := nf.rec.MustGet("funcType").(uuid.UUID)
	ev, err := tx.Schema().GetEnumValueByID(ID(evID))
	if err != nil {
		panic(err)
	}
	return ev
}

func (nf nativeFunction) Call(ctx context.Context, args []interface{}) (interface{}, error) {
	f, ok := nf.nr.fMap[nf.ID()]
	if !ok {
		return nil, fmt.Errorf("func %v not found\n", nf.rec)
	}
	return f(ctx, args)
}
