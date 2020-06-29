package rpc

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
)

var RPCModel = db.MakeModel(
	db.MakeID("29209517-1c39-4be9-9808-e1ed8e40c566"),
	"rpc",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("6ec81a63-0406-4d13-aacf-070a26c2adbc"),
			"name",
			db.String,
		),
	},
	[]db.RelationshipL{RPCCode},
	[]db.ConcreteInterfaceL{},
)

var RPCCode = db.MakeConcreteRelationship(
	db.MakeID("9221119b-495a-449c-b2b3-2c6610f89d7b"),
	"code",
	false,
	db.FunctionInterface,
)

var reactFormRPC = starlark.MakeStarlarkFunction(
	db.MakeID("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
	"reactForm",
	db.RPC,
	`def getDatatypes(name):
    rec = FindOne("model", Eq("name", name))
    # find all attributes associated with the model
    attrs = FindMany("attribute", EqFK("model", rec.ID()))
    out = {}
    for attr in attrs:
        n = str(attr.Get("name")).title()
        dt = FindOne("datatype", EqID(attr.GetFK("datatype")))
        t = str(dt.Get("storedAs"))
        ev = dt.GetFK("enumValues")
        out[n] = {"type" : t, "title" : n, "enum" : ev}
    return out

def makeResponse(input):
    return {
       "title"       : input,
       "description" : "Use form to add new model.",
       "type"        : "object",
       "properties"  : getDatatypes(input)
    }
    
input = args["model"]
result(makeResponse(input))`,
)
