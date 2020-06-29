const CoreApi = {
	datatype:
	{
		create: {},
		findOne: {},
		findMany: {},
		update: {},
		updateMany: {},
	},
	attribute:
	{
		create: {},
		findOne: {},
		findMany: {},
		update: {},
		updateMany: {},
	},
	model:
	{
		create: {},
		findOne: {},
		findMany: {},
		update: {},
		updateMany: {},
	},
	code:
	{
		create: {},
		findOne: {},
		findMany: {},
		update: {},
		updateMany: {},
	},
	rpc:
	{
		create: {},
		findOne: {},
		findMany: {},
		update: {},
		updateMany: {},
	}
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
		this["log"] = async (params) => {
			const res = await fetch(basePath + "log.scan", {
				method: "POST",
				body: JSON.stringify(params)
			});
			const responseBody = await res.json();
			return responseBody.data;
		}
	}
}

var client = new HttpRpcClient("https://localhost:8080/", CoreApi);

module.exports = client;
