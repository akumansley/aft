# Oh we really really need to make this secure
# BIG SCARY COMMENTS
# MASSIVE NEED FOR PERMISSIONS HERE
def main(args):
    out, ran = exec(args["data"], "")
    if not ran:
        return "Starlark: " + out.strip(":")
    return out