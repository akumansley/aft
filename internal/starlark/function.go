package starlark

import (
	"context"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

// Model

var StarlarkFunctionModel = db.MakeModel(
	db.MakeID("c8a17195-b784-4a68-85f4-b4edbfa43174"),
	"starlarkFunction",
	[]db.AttributeL{
		sfName,
		sfCode,
		sfArity,
		sfType,
	},
	[]db.RelationshipL{FunctionRole, StarlarkFunctionModule, ExecutableBy},
	[]db.ConcreteInterfaceL{db.FunctionInterface},
)

var FunctionRole = db.MakeConcreteRelationship(
	db.MakeID("58f08fcb-13ac-43c1-a8db-e1d46114da1b"),
	"role",
	false,
	auth.RoleModel,
)

var ExecutableBy = db.MakeReverseRelationship(
	db.MakeID("89e9f265-bd86-427c-8923-7c09cc7663db"),
	"executableBy",
	auth.ExecutableFunctions,
)

var StarlarkFunctionModule = db.MakeConcreteRelationship(
	db.MakeID("00a50b84-f357-4450-b21b-d0776b97c2c8"),
	"module",
	false,
	db.ModuleModel,
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

var sfType = db.MakeConcreteAttribute(
	db.MakeID("2060e036-2fe3-42ee-a61d-cc00ba3df042"),
	"funcType",
	db.FuncType,
)

func MakeStarlarkFunction(id db.ID, name string, arity int, funcType db.EnumValueL, code string) StarlarkFunctionL {
	return StarlarkFunctionL{
		ID_: id, Name_: name, Code: code, Arity: arity, FuncType: funcType, Role: nil,
	}
}

func MakeStarlarkFunctionWithRole(id db.ID, name string, arity int, funcType db.EnumValueL, code string, role auth.RoleL) StarlarkFunctionL {
	return StarlarkFunctionL{
		id, name, code, arity, funcType, &role,
	}
}

// Literal

type StarlarkFunctionL struct {
	ID_      db.ID         `record:"id"`
	Name_    string        `record:"name"`
	Code     string        `record:"code"`
	Arity    int           `record:"arity"`
	FuncType db.EnumValueL `record:"funcType"`
	Role     *auth.RoleL
}

func (lit StarlarkFunctionL) ID() db.ID {
	return lit.ID_
}

func (lit StarlarkFunctionL) InterfaceID() db.ID {
	return StarlarkFunctionModel.ID()
}

func (lit StarlarkFunctionL) InterfaceName() string {
	return StarlarkFunctionModel.Name_
}

func (lit StarlarkFunctionL) Name() string {
	return lit.Name_
}

func (lit StarlarkFunctionL) Load(tx db.Tx) db.Function {
	f, err := tx.Schema().GetFunctionByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return f
}

func (lit StarlarkFunctionL) Call(ctx context.Context, args []interface{}) (interface{}, error) {
	sr := NewStarlarkRuntime()
	return sr.Execute(ctx, lit.Code, args)
}

func (lit StarlarkFunctionL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit)
	recs = append(recs, rec)
	if lit.Role != nil {
		links = append(links, db.Link{From: lit, To: lit.Role, Rel: FunctionRole})
	}
	return
}

// Runtime

type starlarkFunction struct {
	rec db.Record
	sr  *StarlarkRuntime
}

func (s *starlarkFunction) ID() db.ID {
	return s.rec.ID()
}

func (s *starlarkFunction) Name() string {
	return s.rec.MustGet("name").(string)
}

func (s *starlarkFunction) Code() string {
	return s.rec.MustGet("code").(string)
}

func (s *starlarkFunction) Arity() int {
	a := s.rec.MustGet("arity").(int)
	return a
}

func (s *starlarkFunction) FuncType(tx db.Tx) db.EnumValue {
	evID := s.rec.MustGet("funcType").(uuid.UUID)
	ev, err := tx.Schema().GetEnumValueByID(db.ID(evID))
	if err != nil {
		panic(err)
	}
	return ev
}

func (s *starlarkFunction) Call(ctx context.Context, args []interface{}) (interface{}, error) {
	return s.sr.Execute(ctx, s.Code(), args)
}
