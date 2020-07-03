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
	Code: `def process(name):
    model = FindOne("model", Eq("name", name))
    attrs = FindMany("attribute", EqFK("model", model.ID()))
    schema = {}
    uiSchema = {}
    for attr in attrs:
        name = attr.Get("name")
        dt = FindOne("datatype", EqID(attr.GetFK("datatype")))
        if dt.Get("enum") == False:
            schema[name] = regular(name, dt.Get("storedAs"), dt.Get("name"))
            u = ui(dt.Get("name"))
            if u != None:
                uiSchema[name] = u
        else:
            schema[name] = enum(name, dt.ID(), dt.Get("name"))
    return {"schema" : schema, "uiSchema" : uiSchema}

def regular(fieldName, storedAs, datatype):
    type = FindOne("enumValue", Eq("id",  storedAs)).Get("name")
    if   type == "bool":
          type = "boolean"
    elif type == "int":
         type = "integer"
    elif type == "float":
         type = "number"
    else:
        type  = "string"
    return {
        "type" : type, 
        "title" : fieldName, 
        "datatype" : datatype
    }

def enum(fieldName, id, datatype):
    evs = FindMany("enumValue", EqFK("datatype", id))
    evn = []
    evi = []
    for ev in evs:
        evn.append(ev.Get("name"))
        evi.append(ev.Get("id"))
    return {
        "type" : "string", 
        "title" : fieldName, 
        "enum": evi, 
        "enumNames": evn, 
        "datatype" : datatype
    }

def ui(type):
    if type == "emailAddress":
        return {
            "ui:options": { "inputType": "email"}
        }
    elif type == "url":
        return {"ui:widget" : "uri", "ui:placeholder": "http://"}
    elif type == "bool":
        return {"ui:widget" : "select"}
    elif type == "password":
        return {"ui:widget" : "password"}
    elif type == "longText":
        return {"ui:widget" : "textarea"}
    elif type == "phone" :
        return {
            "ui:options": { "inputType": "tel"}
        }
    return None

out = process(args["model"])
result({
       "schema" : {
       "type"        : "object",
       "properties"  : out["schema"]
    },
       "uiSchema" : out["uiSchema"]})`,
}

var validateRPC = db.Code{
	ID:                db.MakeID("d7633de5-9fa2-4409-a1b2-db96a59be52b"),
	Name:              "validate",
	Runtime:           db.Starlark,
	FunctionSignature: db.RPC,
	Code: `def main(args):
    properties = args["schema"]["properties"]
    data = args["data"]
    errors = {}
    for name in properties:
        x = FindOne("datatype", Eq("name", properties[name]["datatype"]))
        if x.Get("enum") == False:
            y = FindOne("code", EqID(x.GetFK("validator")))
        else:
            y = FindOne("code", Eq("name", "uuid"))
        inp = ""
        if name in data:
            inp = str(data[name])
        out, ran = Exec(y, inp)
        #If there is an error from a validator
        if ran == False:
            errors[name] = {"__errors" : [out]}
    return errors

result(main(args))`,
}

var replRPC = db.Code{
	ID:                db.MakeID("591bc8f7-543b-4fa9-bdf7-8948c79cdd26"),
	Name:              "repl",
	Runtime:           db.Starlark,
	FunctionSignature: db.RPC,
	Code: `def repl(args):
    out, ran = Exec(args["data"], "")
    if ran == False:
        return "Starlark: " + out.strip(":")
    return out

result(repl(args))`,
}

var parseRPC = db.Code{
	ID:                db.MakeID("232d7ad5-357b-43fb-a707-a0a6ba190e7c"),
	Name:              "parse",
	Runtime:           db.Starlark,
	FunctionSignature: db.RPC,
	Code: `# Oh we really really need to make this secure
# BIG SCARY COMMENTS
# MASSIVE NEED FOR PERMISSIONS HERE
def main(args):
    msg, parsed = Parse(args["data"])
    return {"error" : msg, "parsed" : parsed}

result(main(args))`,
}
