---
id: access
title: Access
---

Open up Aft and navigate to the **Access** section.

Create a new role by pressing **Add Role**, and name it `user`. Hit **Save**.

You should see an **Add** button under Grants. Click that, and select `Todo` from the dropdown.

This grants access to all Todos to the user role, which is too broad.

Click **Detail** to expand the policy editor. You should see three sections, labeled Read Create and Update.

Under each, add the following policy.

```js
{
	"user": {"id": "$USER_ID"}
}
```

This indicates that a user can create, read or update only their own Todos.

## Roles

By default, a user has no roles associated with them, so let's connect our test user with the role we just created.

Start by going to the **Terminal** and run the following `update` statement.

```python
def main(aft):
    return aft.api.update("user", {
    	"where": {"email": "user@example.com"},
    	"data": {"role": {"connect": {"name":"user"}}}
    	})
```

Restart Aft as before, but this time with access controls enabled.

```bash
aft -db ./tutorial.dbl -authed=true -serve_dir=client
```

If you open up Aft, you'll notice that no data is displayed, since we don't have access rights to read the schema or other Aft resources. But test out the Todo app, and you should see that everything works!

**Congrats on building your first app with Aft!**

In the last section, we'll review what we've gone over in this tutorial.