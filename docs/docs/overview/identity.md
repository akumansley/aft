---
id: identity
title: Identity
---

Aft has a customizeable login system.

Out of the box, Aft has two [RPCs](rpcs): `login` and `signup`.

The code for each is very short. Here's login:

```python
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
```

There are two methods on the `aft.auth` object: 

1. `authenticateAs(user_id)`, which will generate an authentication token and inject it into the current request, as well as set it in a cookie
2. `user()`, which will return the currently authenticated user object if there is one

The login and signup RPC are just regular code, that you may rewrite to provide any additional logic that's appropriate for your application.

