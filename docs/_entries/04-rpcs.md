---
title: RPCs
order: 4
---

Aft's API can cover most ordinary client needs. But sometimes you just need an escape hatch; for that, Aft includes a scriptable RPC system.

![Screenshot of the rpc page](/aft/img/rpc.png)

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
