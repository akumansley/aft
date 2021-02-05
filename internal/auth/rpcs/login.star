loginUnsuccessful = error(code="login-error", message="login unsuccessful")

def main(args):
    user = findOne("user", {"where": {"email": args["email"]}})
    if not user:
        return loginUnsuccessful
    if func.checkPassword(args["password"], user.id, user.password):
        func.authenticateAs(user.id)
        return user
    else:
        return loginUnsuccessful