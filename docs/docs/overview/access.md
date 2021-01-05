---
id: access
title: Access
---

![Screenshot of the roles page](/img/roles.png)

Aft's powerful API is paired with an equally powerful set of access controls. Aft is closed by default—users with a given role can only access data explicitly granted by a policy.

Policies are expressed on a per-interface and per-operation basis as "where" clauses as used in findMany queries. For example, to restrict access to just users named Andrew, one might have the following read-policy.

```json
{
	"name":"Andrew"
}
```

They're also able to perform template string substitution to restrict on the basis of the current user ID. So to restrict a user to access only their own user data:


```json
{
	"id":"$USER_ID"
}
```

Or for a model that had a relationship to user:

```json
{
	"user": {"id":"$USER_ID"}
}
```

Connections and disconnections are allowed if and only if update is allowed on both records. Similarly, a user must retain update rights to any record that they update after the update is applied.
