const basePath = "https://localhost:8080/";
const methods = ["create", "findOne", "findMany", "update", "updateMany"];

var client = {
  api : new Proxy({}, {
    get: function(target, resource) {
        var out = {};
        methods.forEach((method) => {
          out[method] = async (params) => {
          const res = await fetch(basePath + "api/" + resource + "." + method, {
            method: "POST",
            body: JSON.stringify(params),
          });
          const responseBody = await res.json();
          return responseBody.data;
        }
      });
      return out;
    }
  }),
  views: {
    rpc : new Proxy({}, {
      get: function(target, resource) {
        return async (params) => {
          const res = await fetch(basePath + "views/rpc/" + resource, {
            method: "POST",
            body: JSON.stringify(params),
          });
          const responseBody = await res.json();
          return responseBody.data;
        }
      }
    }),
    repl : async (params) => {
	    const res = await fetch(basePath + "views/repl", {
	      method: "POST",
	      body: JSON.stringify(params)
	    });
	    const responseBody = await res.json();
	    return responseBody;
    }
  },
  log : async (params) => {
	  const res = await fetch(basePath + "log.scan", {
	    method: "POST",
	    body: JSON.stringify(params)
	  });
	  const responseBody = await res.json();
	  return responseBody.data;
  }
};

module.exports = client;
