package rpc

import (
	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
)

var RPCModel = db.MakeModel(
	db.MakeID("25de75bb-2b50-4ec2-87aa-3f1ac74db04a"),
	"rpc",
	[]db.AttributeL{},
	[]db.RelationshipL{RPCFunction, RPCRole},
	[]db.ConcreteInterfaceL{},
)

var RPCFunction = db.MakeConcreteRelationship(
	db.MakeID("a8ffebd7-1805-4714-8d49-7562d30a9c67"),
	"function",
	false,
	db.FunctionInterface,
)

var RPCRole = db.MakeConcreteRelationship(
	db.MakeID("c373f58a-8599-4e57-9a38-2ec1f121f1ef"),
	"role",
	false,
	auth.RoleModel,
)

type RPCL struct {
	ID_      db.ID `record:"id"`
	Function db.FunctionL
	Role     *auth.RoleL
}

func (lit RPCL) ID() db.ID {
	return lit.ID_
}

func (lit RPCL) InterfaceID() db.ID {
	return RPCModel.ID()
}

func (lit RPCL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit)
	f := lit.Function
	links = append(links, db.Link{From: rec.ID(), To: f.ID(), Rel: RPCFunction})
	if lit.Role != nil {
		links = append(links, db.Link{From: rec.ID(), To: lit.Role.ID(), Rel: RPCRole})
	}
	recs = append(recs, rec)
	return
}

func MakeRPC(id db.ID, f db.FunctionL, role *auth.RoleL) RPCL {
	return RPCL{
		ID_:      id,
		Function: f,
		Role:     role,
	}
}
