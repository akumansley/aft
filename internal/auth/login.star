loginUnsuccessful = {"code": "login-error", "message": "login unsuccessful"}

def main(args):
    u = aft.api.findOne("user", {"where": {"email": args.email}}})
    if not u:
        return loginUnsuccessful
    if u.password == args["password"]:
        tok = aft.auth.issueToken(u.id)
        aft.request.setCookie("tok", tok)
        return {
            "data": user,
        }
    else:
        return loginUnsuccessful