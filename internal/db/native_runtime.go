package db

import (
	"github.com/google/uuid"
)

type Func func(interface{}) (interface{}, error)

type NativeRuntime struct {
	fMap map[ID]Func
	db   DB
}

func NewNativeRuntime() *NativeRuntime {
	return &NativeRuntime{
		fMap: make(map[ID]Func),
	}
}

func (nr *NativeRuntime) ProvideModel() Model {
	return NativeFunctionModel
}

func (nr *NativeRuntime) Load(tx Tx, rec Record) Function {
	return nativeFunction{rec, nr, tx}
}

func (nr *NativeRuntime) Save(lit nBox) {
	f := lit.NativeFunctionL.Function
	tx := nr.db.NewRWTx()
	rec, _ := MarshalRecord(lit, NativeFunctionModel)
	tx.Insert(rec)
	nr.fMap[lit.ID()] = f
	tx.Commit()
}

func (nr *NativeRuntime) Registered(db DB) {
	nr.db = db
}

var NativeFunctionModel = ModelL{
	ID:   MakeID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("c47bcd30-01ea-467f-ad02-114342070241"),
			Datatype: String,
		}.AsAttribute(),
		ConcreteAttributeL{
			Name:     "functionSignature",
			ID:       MakeID("ba29d820-ae50-4424-b807-1a1dbd8d2f4b"),
			Datatype: FunctionSignature,
		}.AsAttribute(),
	},
}.AsModel()

type nativeFunction struct {
	rec Record
	nr  *NativeRuntime
	tx  Tx
}

func (nf nativeFunction) ID() ID {
	return nf.rec.ID()
}

func (nf nativeFunction) Name() string {
	return nf.rec.MustGet("name").(string)
}

func (nf nativeFunction) FunctionSignature() EnumValue {
	fs := nf.rec.MustGet("functionSignature").(uuid.UUID)
	ev, _ := nf.tx.Schema().GetEnumValueByID(ID(fs))
	return ev
}

func (nf nativeFunction) Call(args interface{}) (interface{}, error) {
	panic("Not Implemented")
}
