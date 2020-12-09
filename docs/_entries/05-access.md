---
title: Access
order: 5
---

![Screenshot of the roles page](/aft/img/roles.png)

Access controls in Aft play a key role. Becuase of Aft's rich client API, unrestricted API access could allow an end-user to corrupt the application code.

Aft adopts the principle of "closed by default". Users with a given role can only access data explicitly granted by a policy.

Policies are expressed on a per-interface and per-operation basis as "where" clauses as used in findMany queries. They're also able to perform template string substitution to restrict on the basis of the current user ID.

They are also specified separately for read, create and update. Connections and disconnections are allowed if and only if update is allowed on both records.

Similarly, a user must retain update rights to any record that they update after the update.
