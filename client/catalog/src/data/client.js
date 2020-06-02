const CoreApi = {
	log: {
		scan: {},
	},
	model: 
	{
		create: {},
		findOne: {},
		findMany: {},
	},
	datatype:
	{
		create: {},
		findOne: {},
		findMany: {},
	}
}

class HttpRpcClient {
	constructor(basePath, apiSpec) {
		for (let [resource, methods] of Object.entries(apiSpec)) {
			this[resource] = {};
			for (let [method, _] of Object.entries(methods)) {
				this[resource][method] = async (params) => {
					const res = await fetch(basePath + resource + "." + method, {
						method: "POST",
						body: JSON.stringify(params)
					});
					const responseBody = await res.json();
					return responseBody.data;
				}
			}
		}
	}
}

var client = new HttpRpcClient("https://localhost:8080/api/", CoreApi);

module.exports = client;
