def main(args):
    properties = args["schema"]["properties"]
    data = args["data"]
    errors = {}
    for name in properties:
        if name not in data:
            continue
        fn = aft.function.uuid
        if aft.getFunction(name) != None:
             fn = aft.getFunction(name)
        out, err = fn(data)
        if err != None:
            errors[name] = {"__errors" : [err]}
    return errors