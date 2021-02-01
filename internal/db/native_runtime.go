package db

type Func func([]interface{}) (interface{}, error)

// Model

var NativeFunctionModel = MakeModel(
	MakeID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	"nativeFunction",
	[]AttributeL{
		nfName,
		nfArity,
	},
	[]RelationshipL{},
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

func MakeNativeFunction(id ID, name string, arity int, function Func) NativeFunctionL {
	return NativeFunctionL{
		id, name, arity, function,
	}
}

// Literal

type NativeFunctionL struct {
	ID_      ID     `record:"id"`
	Name_    string `record:"name"`
	Arity_   int    `record:"arity"`
	Function Func
}

func (lit NativeFunctionL) ID() ID {
	return lit.ID_
}

func (lit NativeFunctionL) InterfaceID() ID {
	return NativeFunctionModel.ID()
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

func (nr *NativeRuntime) Save(b *Builder, lit NativeFunctionL) {
	f := lit.Function
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

func (nf nativeFunction) Call(args []interface{}) (interface{}, error) {
	return nf.nr.fMap[nf.ID()](args)
}
