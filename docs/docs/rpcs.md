---
id: rpcs
title: RPCs
---


Aft's API can cover most ordinary client needs. But sometimes you just need an escape hatch; for that, Aft includes a scriptable RPC system.

![Screenshot of the rpc page](/img/rpc.png)

RPCs can be written in Starlark. They are passed two arguments; first, a handle to a Starlark version of the Aft API and second, a dictionary that is a decoded version of a JSON object passed from the client.

The RPCs are exposed in the following URL format:

```
https://$BASE_URL/api/rpc.$RPC_NAME
```

The RPC endpoint accepts a JSON object with a single key, "args":

```
{
	"args": {
		"foo": "bar"
	}
}
```

## Starlark

![Screenshot of the rpc edit page](/img/rpcedit.png)

To write an RPC in Starlark, write a script with a function, "main," of two arguments: `aft`, a handle to the API and authentication methods, and `data`, a single json object sent by the client.

Aft's API methods can be accessed by calling them on the `api` object like so:

```python
def main(aft, data):
	user = aft.api.findOne("users", {"where": {"name": "Andrew"}})
	return user.name  # returns "Andrew"
```

Modifying records returned by the `aft.api` methods will not have an effect on the datastore. To mutate the datastore, use the mutation api calls. 

```python
def main(aft, data):
	user = aft.api.findOne("users", {"where": {"name": "Andrew"}})
	aft.api.update("users", {"where": {"id": user.id}, "data": {"name": "Werdna"}})
```
