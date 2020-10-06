package rpc

import (
    "awans.org/aft/internal/db"
    "awans.org/aft/internal/starlark"
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
