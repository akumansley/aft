def main(args)
    out = process(args["model"])
    return {
       "schema" : {
           "type"        : "object",
           "properties"  : out["schema"]
		},
       "uiSchema" : out["uiSchema"]
    }
    
def process(name):
    model = aft.api.findOne("model", {
    	"where" : {"name" : name}, 
    	"include" : {"attributes" : {"include" : {"datatype" : True}}}
    })
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