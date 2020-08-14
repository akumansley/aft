def main(args):
    msg, parsed = parse(args["data"])
    if parsed:
        return {
            "message" : msg,
            "parsed" : parsed
        }
    parts = msg.split(":",3)
    return {
        "message" : parts[3].strip().title(), 
        "parsed" : False,
        "line" : int(parts[1]),
        "start" : int(parts[2]),
        "raw" : msg
    }