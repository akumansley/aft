const CoreApi = {
	log: {
		scan: {},
	},
	datatype:
	{
		create: {},
		findOne: {},
		findMany: {},
	},
	model:
	{
		create: {},
		findOne: {},
		findMany: {},
	},
	code:
	{
		create: {},
		findOne: {},
		findMany: {},
	},
	rpc:
	{
		create: {},
		findOne: {},
		findMany: {},
	},
}

class HttpRpcClient {
	constructor(basePath, apiSpec) {
		for (let [resource, methods] of Object.entries(apiSpec)) {
			this[resource] = {};
			for (let [method, _] of Object.entries(methods)) {
				this[resource][method] = async (params) => {
					const res = await fetch(basePath + "api/" + resource + "." + method, {
						method: "POST",
						body: JSON.stringify(params)
					});
					const responseBody = await res.json();
					return responseBody.data;
				}
			}
		}
		this["repl"] = async (params) => {
			const res = await fetch(basePath + "views/repl", {
				method: "POST",
				body: JSON.stringify(params)
			});
			const responseBody = await res.json();
			return responseBody;
		}
	}
}

var client = new HttpRpcClient("https://localhost:8080/", CoreApi);

module.exports = client;
