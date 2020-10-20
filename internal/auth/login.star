loginUnsuccessful = {"code": "login-error", "message": "login unsuccessful"}

def main(aft, args):
    user = aft.api.findOne("user", {"where": {"email": args["email"]}})
    if not user:
        return loginUnsuccessful
    if user.password == args["password"]:
        aft.auth.authenticateAs(user.id)
        return user
    else:
        return loginUnsuccessful