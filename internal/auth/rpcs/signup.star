signupUnsuccessful = error(code="signup-error", message="signup unsuccessful")

def main(args):
	user = create("user", {"data": {
			"email": args["email"], 
			"password": args["password"],
			"role": {
				"connect": {"name": "user"}
			}
		}})

	if not user:
		return signupUnsuccessful
	func.authenticateAs(user.id)
	return user
