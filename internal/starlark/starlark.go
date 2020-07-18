package starlark

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/chasehensel/starlight/convert"
	"github.com/google/uuid"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

// Model

var StarlarkFunctionModel = db.MakeModel(
	db.MakeID("c8a17195-b784-4a68-85f4-b4edbfa43174"),
	"starlarkFunction",
	[]db.AttributeL{
		sfName,
		sfCode,
		sfFuncSig,
	},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{db.FunctionInterface},
)

var sfCode = db.MakeConcreteAttribute(
	db.MakeID("8e6538a9-f64e-4b48-986c-55924bd1da2d"),
	"code",
	db.String,
)

var sfName = db.MakeConcreteAttribute(
	db.MakeID("70b6e7e8-e47c-488c-81b3-e76c0eac0891"),
	"name",
	db.String,
)

var sfFuncSig = db.MakeConcreteAttribute(
	db.MakeID("1530c2cf-b61b-4d20-a130-77bba8d203b1"),
	"functionSignature",
	db.FunctionSignature,
)

// Loader

func NewStarlarkRuntime() *StarlarkRuntime {
	return &StarlarkRuntime{}
}

type StarlarkRuntime struct {
}

//configure starlark
func init() {
	resolve.AllowNestedDef = true // allow def statements within function bodies
	resolve.AllowLambda = true    // allow lambda expressions
	resolve.AllowFloat = true     // allow floating point literals, the 'float' built-in, and x / y
	resolve.AllowSet = true       // allow the 'set' built-in
	resolve.AllowRecursion = true // allow while statements and recursive functions
}

func (sr *StarlarkRuntime) Execute(code string, functionSignature db.EnumValue, input interface{}, env map[string]interface{}) (interface{}, error) {
	i, err := convert.ToValue(input)
	if err != nil {
		return nil, err
	}
	c := &call{}
	c.Env = env
	globals, err := CreateEnv(i, c)
	if err != nil {
		return nil, err
	}

	if functionSignature.ID() == db.FromJSON.ID() {
		code = fmt.Sprintf("%s\n\nresult(validator(args))", code)
	}
	// Run the starlark interpreter!
	_, err = starlark.ExecFile(&starlark.Thread{Load: nil}, "", []byte(code), globals)

	if err != nil {
		if evale, ok := err.(*starlark.EvalError); ok {
			return nil, fmt.Errorf("\n%s", evale.Backtrace())
		}
		return nil, err
	}
	if c.err != nil {
		return nil, fmt.Errorf("Raised: %s", c.err)
	}
	out := recursiveFromValue(c.result)
	return out, nil
}

func (sr *StarlarkRuntime) ProvideModel() db.ModelL {
	return StarlarkFunctionModel
}

func (sr *StarlarkRuntime) Load(tx db.Tx, rec db.Record) db.Function {
	return &starlarkFunction{rec, sr, tx}
}

func MakeStarlarkFunction(id db.ID, name string, functionSignature db.EnumValue, code string) StarlarkFunctionL {
	return StarlarkFunctionL{
		id, name, code, functionSignature,
	}
}

type StarlarkFunctionL struct {
	ID_                db.ID        `record:"id"`
	Name_              string       `record:"name"`
	Code               string       `record:"code"`
	FunctionSignature_ db.EnumValue `record:"functionSignature"`
}

func (lit StarlarkFunctionL) ID() db.ID {
	return lit.ID_
}

func (lit StarlarkFunctionL) Name() string {
	return lit.Name_
}

func (lit StarlarkFunctionL) FunctionSignature() db.EnumValue {
	return lit.FunctionSignature_
}

func (lit StarlarkFunctionL) Call(args interface{}) (interface{}, error) {
	sr := NewStarlarkRuntime()
	return sr.Execute(lit.Code, lit.FunctionSignature(), args, nil)
}

func (lit StarlarkFunctionL) CallWithEnv(args interface{}, env map[string]interface{}) (interface{}, error) {
	sr := NewStarlarkRuntime()
	return sr.Execute(lit.Code, lit.FunctionSignature(), args, env)
}

func (lit StarlarkFunctionL) MarshalDB() (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(lit, StarlarkFunctionModel)
	recs = append(recs, rec)
	return
}

type starlarkFunction struct {
	rec db.Record
	sr  *StarlarkRuntime
	tx  db.Tx
}

func (s *starlarkFunction) ID() db.ID {
	return s.rec.ID()
}

func (s *starlarkFunction) Name() string {
	return sfName.MustGet(s.rec).(string)
}

func (s *starlarkFunction) Code() string {
	return sfCode.MustGet(s.rec).(string)
}

func (s *starlarkFunction) FunctionSignature() db.EnumValue {
	// TODO think of a better way to handle reading enums out of structs
	u := sfFuncSig.MustGet(s.rec).(uuid.UUID)
	ev, _ := s.tx.Schema().GetEnumValueByID(db.ID(u))
	return ev
}

func (s *starlarkFunction) Call(input interface{}) (interface{}, error) {
	return s.sr.Execute(s.Code(), s.FunctionSignature(), input, nil)
}
