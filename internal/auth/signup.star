signupUnsuccessful = {"code": "signup-error", "message": "signup unsuccessful"}

def main(aft, args):
	user = aft.api.create("user", {"data": {
			"email": args["email"], 
			"password": args["password"],
			"role": {
				"connect": {"name": "user"}
			}
		}})

	if not user:
		return signupUnsuccessful
	aft.auth.authenticateAs(user.id)
	return user
