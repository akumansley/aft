def main(args):
    properties = args["schema"]["properties"]
    data = args["data"]
    errors = {}
    for name in properties:
        if name not in data:
            continue
        x findOne("datatype", {"where" : {"name" : properties[name]["datatype"]}, "include" : {"validator" : True}})
        if x.enum:
            y = findOne("code", {"where" : {"name": "uuid"}}).code
        else:
            y = x.validator.code
        out, success = getattr(func, name)()
        # If there is an error from a validator
        if success == False:
            errors[name] = {"__errors" : [out]}
    return errors
