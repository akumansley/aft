package starlark

import (
	"awans.org/aft/internal/db"
)

// Model

var StarlarkFunctionModel = db.MakeModel(
	db.MakeID("c8a17195-b784-4a68-85f4-b4edbfa43174"),
	"starlarkFunction",
	[]db.AttributeL{
		sfName,
		sfCode,
		sfArity,
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

var sfArity = db.MakeConcreteAttribute(
	db.MakeID("9aca21f6-3fc6-4b17-a4e7-1674bd6a7593"),
	"arity",
	db.Int,
)

func MakeStarlarkFunction(id db.ID, name string, arity int, code string) StarlarkFunctionL {
	return StarlarkFunctionL{
		id, name, code, arity,
	}
}

// Literal

type StarlarkFunctionL struct {
	ID_    db.ID  `record:"id"`
	Name_  string `record:"name"`
	Code   string `record:"code"`
	Arity_ int    `record:"arity"`
}

func (lit StarlarkFunctionL) ID() db.ID {
	return lit.ID_
}

func (lit StarlarkFunctionL) Name() string {
	return lit.Name_
}

func (lit StarlarkFunctionL) Arity() int {
	return lit.Arity_
}

func (lit StarlarkFunctionL) Call(args []interface{}) (interface{}, error) {
	sr := NewStarlarkRuntime()
	return sr.Execute(lit.Code, args)
}

func (lit StarlarkFunctionL) MarshalDB() (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(lit, StarlarkFunctionModel)
	recs = append(recs, rec)
	return
}

// Runtime

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

func (s *starlarkFunction) Arity() int {
	a := sfArity.MustGet(s.rec).(int)
	return a
}

func (s *starlarkFunction) Call(args []interface{}) (interface{}, error) {
	return s.sr.Execute(s.Code(), args)
}