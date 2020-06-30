import { isObject } from "./util";

const basePath = "https://e329e49c8232.ngrok.io/";
const methods = ["create", "findOne", "findMany", "update", "updateMany"];

var client = {
  api: new Proxy(
    {},
    {
      get: function(target, resource) {
        var out = {};
        methods.forEach(method => {
          out[method] = async params => {
            const res = await fetch(
              basePath + "api/" + resource + "." + method,
              {
                method: "POST",
                body: JSON.stringify(params)
              }
            );
            const responseBody = await res.json();
            if ("code" in responseBody) {
              var e = new Error(responseBody.message);
              e.name = responseBody.code;
              return Promise.reject(e);
            }
            return responseBody.data;
          };
        });
        return out;
      }
    }
  ),
  views: {
    rpc: new Proxy(
      {},
      {
        get: function(target, resource) {
          return async params => {
            const res = await fetch(basePath + "views/rpc/" + resource, {
              method: "POST",
              body: JSON.stringify(params)
            });
            const responseBody = await res.json();
            if ("code" in responseBody) {
              var e = new Error(responseBody.message);
              e.name = responseBody.code;
              return Promise.reject(e);
            }
            return responseBody.data;
          };
        }
      }
    ),
    repl: async params => {
      const res = await fetch(basePath + "views/repl", {
        method: "POST",
        body: JSON.stringify(params)
      });
      const responseBody = await res.json();
      if (isObject(responseBody)) {
        var e = new Error(responseBody.message);
        e.name = responseBody.code;
        return Promise.reject(e);
      }
      return responseBody;
    }
  },
  log: async params => {
    const res = await fetch(basePath + "log.scan", {
      method: "POST",
      body: JSON.stringify(params)
    });
    const responseBody = await res.json();
    if ("code" in responseBody) {
      var e = new Error(responseBody.message);
      e.name = responseBody.code;
      return Promise.reject(e);
    }
    return responseBody.data;
  }
};

export default client;