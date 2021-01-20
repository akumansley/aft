loginUnsuccessful = error(code="login-error", message="login unsuccessful")

def main(aft, args):
    user = aft.api.findOne("user", {"where": {"email": args["email"]}})
    if not user:
        return loginUnsuccessful
    if aft.auth.checkPassword(args["password"], user.id, user.password):
        aft.auth.authenticateAs(user.id)
        return user
    else:
        return loginUnsuccessful