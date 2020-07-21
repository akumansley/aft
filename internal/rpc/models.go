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
	`def process(name):
    model = FindOne("model", Eq("name", name))
    attrs = FindMany("attribute", EqFK("model", model.ID()))
    schema = {}
    uiSchema = {}
    for attr in attrs:
        name = attr.Get("name")
        dt = FindOne("datatype", EqID(attr.GetFK("datatype")))
        if not dt.Get("enum"):
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
)

var validateFormRPC = starlark.MakeStarlarkFunction(
	db.MakeID("d7633de5-9fa2-4409-a1b2-db96a59be52b"),
	"validateForm",
	db.RPC,
	`def main(args):
    properties = args["schema"]["properties"]
    data = args["data"]
    errors = {}
    for name in properties:
        if name not in data:
            continue
        x = FindOne("datatype", Eq("name", properties[name]["datatype"]))
        if not x.Get("enum"):
            y = FindOne("code", EqID(x.GetFK("validator")))
        else:
            y = FindOne("code", Eq("name", "uuid"))
        out, success = Exec(y, data[name])
        #If there is an error from a validator
        if ran == False:
            errors[name] = {"__errors" : [out]}
    return errors

result(main(args))`,
)

var replRPC = starlark.MakeStarlarkFunction(
	db.MakeID("591bc8f7-543b-4fa9-bdf7-8948c79cdd26"),
	"repl",
	db.RPC,
	`# Oh we really really need to make this secure
        if not success:
            out = out.split("fail: ")
            errors[name] = {"__errors" : [out[-1]]}
    return errors`,
)

var terminalRPC = starlark.MakeStarlarkFunction(
	db.MakeID("180bb262-8835-4ec5-9c2b-3f455615be9a"),
	"terminal",
	db.RPC,
	`# Oh we really really need to make this secure
# BIG SCARY COMMENTS
# MASSIVE NEED FOR PERMISSIONS HERE
def main(args):
    out, ran = Exec(args["data"], "")
    if not ran:
        return "Starlark: " + out.strip(":")
    return out`,
)

var parseRPC = starlark.MakeStarlarkFunction(
	db.MakeID("232d7ad5-357b-43fb-a707-a0a6ba190e7c"),
	"parse",
	db.RPC,
	`def main(args):
    msg, parsed = Parse(args["data"])
    return {"error" : msg, "parsed" : parsed}
    return out`,
)

var lintRPC = starlark.MakeStarlarkFunction(
	db.MakeID("e4be72dc-9462-49f7-bba9-3543cc6bf6c2"),
	"lint",
	db.RPC,
	`def main(args):
    msg, parsed = Parse(args["data"])
    if parsed:
        return {
            "message" : msg,
            "parsed" : parsed
        }
    parts = msg.split(":",3)
    return {
        "message" : parts[3].strip().title(), 
        "parsed" : parsed,
        "line" : int(parts[1]),
        "start" : int(parts[2]),
        "raw" : msg
    }`,
)
