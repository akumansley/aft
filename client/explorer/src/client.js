const basePath = "https://localhost:8080/";
const methods = ["create", "findOne", "findMany", "update", "updateMany", "count"];

async function call(path, params) {
  if(typeof params === 'undefined')  {
    params = {};
  }
  const res = await fetch(basePath + path, {
    method: "POST",
    body: JSON.stringify(params)
  });
  const responseBody = await res.json();
  if ("code" in responseBody) {
    return Promise.reject(responseBody);
  }
  if("data" in responseBody) {
    return responseBody.data;
  }
  if("BatchPayload" in responseBody) {
    return responseBody.BatchPayload;
  }
}

var client = {
  api: new Proxy(
    {},
    {
      get: function(target, resource) {
        var out = {};
        methods.forEach(method => {
          out[method] = params => {
            return call("api/" + resource + "." + method, params);
          };
        });
        return out;
      }
    }
  ),
  rpc: new Proxy(
    {},
    {
      get: function(target, resource) {
        return params => {
          return call("rpc/" + resource, params);
        };
      }
    }
  ),
  log: params => {
    return call("log.scan", params);
  }
};

export default client;
