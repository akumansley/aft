function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) {
    return parts.pop().split(";").shift();
  }
  return ""
};

async function call(path, body) {
	const result = await fetch(path, {
		method: 'POST',
		headers: {'X-CSRF': getCookie('csrf')},
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
