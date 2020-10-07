loginUnsuccessful = {"code": "login-error", "message": "login unsuccessful"}

def main(aft, args):
    u = aft.api.findOne("user", {"where": {"email": args.email}}})
    if not u:
        return loginUnsuccessful
    if u.password == args["password"]:
        aft.auth.authenticateAs(u.id)
        return {
            "data": user,
        }
    else:
        return loginUnsuccessful