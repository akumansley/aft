async function call(path, body) {
	const result = await fetch(path, {
		method: 'POST',
		headers: {'Content-Type': 'application/json'},
		body: JSON.stringify(body || {}),
	})
	const response = await result.json();
	if (response.code) {
		throw new Error(response.message);
	}
	return response.data;
}

const curryProxy = (inner) => {
	return new Proxy({}, {
		get(_, prop)  { 
			return inner(prop) 
		}
	})
}

export default {
	api: curryProxy((interfaceName) => curryProxy((method) => (params) => {
		return call("api/" + interfaceName + "." + method, params)
	})),
	rpc: curryProxy((rpcName) => (args) => {
		return call("rpc/" + rpcName, args)
	}),
}
