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
    model = aft.api.findOne("model", {"where" : {"name" : name}, "include" : {"attributes" : {"include" : {"datatype" : True}}}})
    schema = {}
    uiSchema = {}
    for attr in model.attributes:
        if not attr.datatype.enum:
            schema[attr.name] = regular(attr.name, attr.datatype)
            u = ui(attr.datatype.name)
            if u != None:
                uiSchema[attr.name] = u
        else:
            schema[attr.name] = enum(attr.name, attr.datatype)
    return {"schema" : schema, "uiSchema" : uiSchema}

def regular(fieldName, datatype):
    type = aft.api.findOne("enumValue", {"where" : {"id" : datatype.storedAs}}).name
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
        "datatype" : datatype.name
    }

def enum(fieldName, datatype):
    dt = aft.api.findOne("datatype", {"where" : {"id" : datatype.id}, "include" : {"enumValues" : True}})
    evn = []
    evi = []
    for ev in dt.enumValues:
        evn.append(ev,name)
        evi.append(ev.id)
    return {
        "type" : "string", 
        "title" : fieldName, 
        "enum": evi, 
        "enumNames": evn, 
        "datatype" : datatype.name
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

def main(args)
    out = process(args["model"])
    return {
       "schema" : {
       "type"        : "object",
       "properties"  : out["schema"]
    },
       "uiSchema" : out["uiSchema"]}`,
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
        x = aft.api.findOne("datatype", {"where" : {"name" : properties[name]["datatype"]}, "include" : {"validator" : True}})
        if x.enum:
            y = aft.api.findOne("code", {"where" : {"name": "uuid"}}).code
        else:
            y = x.validator.code
        out, success = exec(y, data[name])
        #If there is an error from a validator
        if ran == False:
            errors[name] = {"__errors" : [out]}
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
    out, ran = exec(args["data"], "")
    if not ran:
        return "Starlark: " + out.strip(":")
    return out`,
)

var parseRPC = starlark.MakeStarlarkFunction(
	db.MakeID("232d7ad5-357b-43fb-a707-a0a6ba190e7c"),
	"parse",
	db.RPC,
	`def main(args):
    msg, parsed = parse(args["data"])
    return {"error" : msg, "parsed" : parsed}
    return out`,
)

var lintRPC = starlark.MakeStarlarkFunction(
	db.MakeID("e4be72dc-9462-49f7-bba9-3543cc6bf6c2"),
	"lint",
	db.RPC,
	`def main(args):
    msg, parsed = parse(args["data"])
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
