package rpc

import (
	"awans.org/aft/internal/db"
)

var RPCModel = db.Model{
	ID:   db.MakeModelID("29209517-1c39-4be9-9808-e1ed8e40c566"),
	Name: "rpc",
	Attributes: []db.Attribute{
		db.Attribute{
			Name:     "name",
			ID:       db.MakeID("6ec81a63-0406-4d13-aacf-070a26c2adbc"),
			Datatype: db.String,
		},
	},
}

var RPCCode = db.Relationship{
	ID:     db.MakeID("9221119b-495a-449c-b2b3-2c6610f89d7b"),
	Source: RPCModel,
	Name:   "code",
	Multi:  false,
	Target: db.CodeModel,
}

var reactFormRPC = db.Code{
	ID:                db.MakeID("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
	Name:              "reactForm",
	Runtime:           db.Starlark,
	FunctionSignature: db.RPC,
	Code: `def getDatatypes(name):
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
}
