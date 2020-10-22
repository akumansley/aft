const apiBase = "/api"
const apiMethods = ["create", "update", "findOne", "findMany"];
const objects = ["gift", "user"];

const rpcBase = "/rpc";
const rpcMethods = ["login", "signup", "me"]

const basePath = "https://localhost:8080";

function getToken() {
	let cookie = document.cookie;
	if (cookie) {
		let tok = cookie.split('; ')
		.find(row => row.startsWith('tok'))
		.split('=')[1];
		return tok;
	}
	return "";
}


async function post(url, params) {
	if(typeof params === 'undefined')  {
		params = {};
	}
	try {
		const res = await fetch(url, {
			method: "POST",
			body: JSON.stringify(params),
			headers: new Headers({
				'Authorization': getToken(),
			}),
			credentials: 'include',
		});
		const responseBody = await res.json();
		if ("code" in responseBody) {
			return Promise.reject(responseBody);
		}
		if("data" in responseBody) {
			return responseBody.data;
		}
	} catch (err) {
		console.log(err);
		throw err;
	}

}

function api(objects, methods) {
	const a = {};
	for (let o of objects) {
		a[o] = {};
		for (let m of methods) {
			a[o][m] = (params) => {
				return post(basePath + apiBase + '/' + o + "." + m, params);
			}
		}
	}
	return a;
}

function rpcs(methods) {
	const v = {};
	for (let m of methods) {
		v[m] = (params) => {
			return post(basePath + rpcBase + '/' + m, {"args": params});
		}
	}
	return v;
}

const client = {
	api: api(objects, apiMethods),
	rpc: rpcs(rpcMethods),
}

export default client;

