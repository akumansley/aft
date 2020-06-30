package rpc

import (
	"awans.org/aft/internal/db"
)

var RPCModel = db.Model{
	ID:     db.MakeModelID("29209517-1c39-4be9-9808-e1ed8e40c566"),
	Name:   "rpc",
	System: true,
	Attributes: []db.Attribute{
		db.Attribute{
			Name:     "name",
			ID:       db.MakeID("6ec81a63-0406-4d13-aacf-070a26c2adbc"),
			Datatype: db.String,
		},
	},
	LeftRelationships: []db.Relationship{
		RPCCode,
	},
}

var RPCCode = db.Relationship{
	ID:           db.MakeID("9221119b-495a-449c-b2b3-2c6610f89d7b"),
	LeftModelID:  db.MakeModelID("29209517-1c39-4be9-9808-e1ed8e40c566"), // rpc
	LeftName:     "code",
	LeftBinding:  db.BelongsTo,
	RightModelID: db.MakeModelID("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "rpc",
	RightBinding: db.HasOne,
}

var reactFormRPC = db.Code{
	ID:                db.MakeID("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
	Name:              "reactForm",
	Runtime:           db.Starlark,
	FunctionSignature: db.RPC,
	Code: `def properties(name):
    model = FindOne("model", Eq("name", name))
    attrs = FindMany("attribute", EqFK("model", model.ID()))
    out = {}
    for attr in attrs:
        name = str(attr.Get("name"))
        dt = FindOne("datatype", EqID(attr.GetFK("datatype")))
        if str(dt.Get("enum")) == "false":
            out[name] = regular(name, dt.Get("storedAs"))
        else:
            out[name] = enum(name, dt.ID())
    return out

def regular(fieldName, storedAs):
    tr = FindOne("enumValue", Eq("id",  storedAs))
    t = str(tr.Get("name"))
    if t == "bool":
        t = "boolean"
    elif t == "int":
        t = "integer"
    elif t == "float":
        t = "number"
    else:
        t = "string"
    return {"type" : t, "title" : fieldName}

def enum(fieldName, id):
    evs = FindMany("enumValue", EqFK("datatype", id))
    evn = []
    evi = []
    for ev in evs:
        evn.append(ev.Get("name"))
        evi.append(ev.Get("id"))
    return {"type" : "string", "title" : fieldName, "enum": evi, "enumNames": evn}

result({
       "type"        : "object",
       "properties"  : properties(args["model"])
    })`,
}
