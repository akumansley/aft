export let cap= (s) => {
	if (!s) {
		return "";
	}
	s = s.replace(/([A-Z]+)/g, " $1").replace(/([A-Z][a-z])/g, " $1");
	return s.charAt(0).toUpperCase() + s.slice(1)
}

export let restrictToIdent= (s) => {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal.toLowerCase();
}

export let getEnumsFromObj = (obj) => {
	var runtime = {};
	var fs = {};
	var storage = {};
	for (var i = 0; i < obj.length; i++) {
		var name = obj[i]["name"];
		if(name == "runtime") {
			var enumValues = obj[i]["enumValues"];
			for (var j = 0; j < enumValues.length; j++) {
				runtime[enumValues[j]["name"]] = enumValues[j];
			}		
		} else if(name == "functionSignature") {
			var enumValues = obj[i]["enumValues"];
			for (var j = 0; j < enumValues.length; j++) {
				fs[enumValues[j]["name"]] = enumValues[j];
			}		
		} else if(name == "storedAs") {
			var enumValues = obj[i]["enumValues"];
			for (var j = 0; j < enumValues.length; j++) {
				storage[enumValues[j]["name"]] = enumValues[j];
			}		
		}
	}
	return {"runtime" : runtime, "fs" : fs, "storage" : storage}
}

export let isObject = s => {
  return typeof s == "object";
};
