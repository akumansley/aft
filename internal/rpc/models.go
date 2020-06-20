package rpc

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

var RPCModel = db.Model{
	ID:   uuid.MustParse("29209517-1c39-4be9-9808-e1ed8e40c566"),
	Name: "rpc",
	Attributes: map[string]db.Attribute{
		"name": db.Attribute{
			ID:       uuid.MustParse("6ec81a63-0406-4d13-aacf-070a26c2adbc"),
			Datatype: db.String,
		},
	},
	LeftRelationships: []db.Relationship{
		RPCCode,
	},
}

var RPCCode = db.Relationship{
	ID:           uuid.MustParse("9221119b-495a-449c-b2b3-2c6610f89d7b"),
	LeftModelID:  uuid.MustParse("29209517-1c39-4be9-9808-e1ed8e40c566"), // rpc
	LeftName:     "code",
	LeftBinding:  db.BelongsTo,
	RightModelID: uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "rpc",
	RightBinding: db.HasOne,
}

var reactFormRPC = db.Code{
	ID:                uuid.MustParse("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
	Name:              "reactForm",
	Runtime:           db.Starlark,
	FunctionSignature: db.RPC,
	Code: `def getDatatypes(name):
    rec = FindOne("model", Eq("name", name))
    # find all attributes associated with the model
    attrs = FindMany("attribute", EqFK("model", rec.ID()))
    out = {}
    for attr in attrs:
        n = str(attr.Get("name"))
        dt = FindOne("datatype", Eq("id", attr.GetFK("datatype")))
        t = str(dt.Get("name"))
        if t == "enum" or t == "int":
            t = "integer"
        if t == "emailAddress" or "url":
            t = "string"
        out[n] = {"type" : t, "title" : n}
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
