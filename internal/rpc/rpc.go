package rpc

import "awans.org/aft/internal/db"

var RPCModel = db.MakeModel(
	db.MakeID("25de75bb-2b50-4ec2-87aa-3f1ac74db04a"),
	"rpc",
	[]db.AttributeL{},
	[]db.RelationshipL{RPCFunction},
	[]db.ConcreteInterfaceL{},
)

var RPCFunction = db.MakeConcreteRelationship(
	db.MakeID("a8ffebd7-1805-4714-8d49-7562d30a9c67"),
	"function",
	false,
	db.FunctionInterface,
)

type RPCL struct {
	ID_      db.ID `record:"id"`
	Function db.FunctionL
}

func (lit RPCL) ID() db.ID {
	return lit.ID_
}

func (lit RPCL) MarshalDB() (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(lit, RPCModel)
	f := lit.Function
	links = append(links, db.Link{From: rec.ID(), To: f.ID(), Rel: RPCFunction})
	recs = append(recs, rec)
	return
}

func MakeRPC(id db.ID, f db.FunctionL) RPCL {
	return RPCL{
		ID_:      id,
		Function: f,
	}
}
