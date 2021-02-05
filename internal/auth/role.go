package auth

import (
	"awans.org/aft/internal/db"
)

var RoleModel = db.MakeModel(
	db.MakeID("bf17994e-7ef1-459f-9b82-069016686081"),
	"role",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("6dc3ec26-3125-4e54-b9b0-f5ccad10c4af"),
			"name",
			db.String,
		),
	},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

func init() {
	RoleModel.Relationships_ = []db.RelationshipL{
		RolePolicy, RoleUsers, ExecutableFunctions,
	}
	PolicyModel.Relationships_ = []db.RelationshipL{
		PolicyRole, PolicyFor,
	}
}

var ExecutableFunctions = db.MakeConcreteRelationship(
	db.MakeID("f548851c-d548-472a-9aaf-6cf8260586a8"),
	"executableFunctions",
	true,
	db.FunctionInterface,
)

var ExecutableBy = db.MakeInterfaceRelationshipWithSource(
	db.MakeID("2ccd7e75-740e-4a36-9991-4dfdba6a5df4"),
	"executableBy",
	true,
	db.FunctionInterface,
	RoleModel,
)

var NativeFunctionExecutableBy = db.MakeReverseRelationshipWithSource(
	db.MakeID("b30827ee-a8b2-4283-8600-4395dc7515e1"),
	"executableBy",
	ExecutableFunctions,
	db.NativeFunctionModel,
)

var RolePolicy = db.MakeConcreteRelationship(
	db.MakeID("fc193452-3c43-4019-b886-d95decc1ce97"),
	"policies",
	true,
	PolicyModel,
)

var FunctionRole = db.MakeInterfaceRelationshipWithSource(
	db.MakeID("11dc7bf0-c30e-4d7f-afc7-fefa412b7583"),
	"role",
	false,
	db.FunctionInterface,
	RoleModel,
)

var NativeFunctionRole = db.MakeConcreteRelationshipWithSource(
	db.MakeID("8b808f27-0e48-4a86-bbb6-224a962fa8d7"),
	"role",
	false,
	db.NativeFunctionModel,
	RoleModel,
)

var RoleUsers = db.MakeReverseRelationship(
	db.MakeID("098dd9f8-1337-44b2-bf8d-277e4aafd725"),
	"users",
	UserRole,
)

var RoleModule = db.MakeConcreteRelationshipWithSource(
	db.MakeID("9a3a58a9-2854-4e22-adfd-1d7352089a0a"),
	"module",
	false,
	RoleModel,
	db.ModuleModel,
)

var ModuleRoles = db.MakeReverseRelationshipWithSource(
	db.MakeID("0bf7990c-985a-4fb7-aad7-c4b459166a46"),
	"roles",
	RoleModule,
	db.ModuleModel,
)

func MakeNativeFunctionWithRole(id db.ID, name string, arity int, funcType db.EnumValueL, function db.Func, role RoleL) NativeFunctionWithRoleL {
	return NativeFunctionWithRoleL{
		function: db.NativeFunctionL{
			id, name, arity, funcType, function,
		},
		role: role,
	}
}

// Literal
type NativeFunctionWithRoleL struct {
	function db.NativeFunctionL
	role     RoleL
}

func (lit NativeFunctionWithRoleL) ID() db.ID {
	return lit.function.ID()
}

func (lit NativeFunctionWithRoleL) InterfaceID() db.ID {
	return lit.function.InterfaceID()
}

func (lit NativeFunctionWithRoleL) InterfaceName() string {
	return lit.function.InterfaceName()
}

func (lit NativeFunctionWithRoleL) Func() db.Func {
	return lit.function.Func()
}

func (lit NativeFunctionWithRoleL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit.function)
	recs = append(recs, rec)
	links = append(links, db.Link{lit, lit.role, NativeFunctionRole})
	return
}

func (lit NativeFunctionWithRoleL) Load(tx db.Tx) db.Function {
	return lit.function.Load(tx)
}

type RoleL struct {
	ID_       db.ID  `record:"id"`
	Name      string `record:"name"`
	Policies  []PolicyL
	Functions []db.FunctionL
}

func (lit RoleL) ID() db.ID {
	return lit.ID_
}

func (lit RoleL) InterfaceID() db.ID {
	return RoleModel.ID()
}

func (lit RoleL) InterfaceName() string {
	return RoleModel.Name_
}

func (lit RoleL) ModuleRelationship() db.ConcreteRelationshipL {
	return RoleModule
}

func (lit RoleL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit)
	for _, p := range lit.Policies {
		links = append(links, db.Link{From: lit, To: p, Rel: RolePolicy})
	}
	for _, f := range lit.Functions {
		links = append(links, db.Link{From: lit, To: f, Rel: ExecutableFunctions})
	}
	recs = append(recs, rec)
	return
}
