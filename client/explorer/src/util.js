export let cap = s => {
  if (!s) {
    return "";
  }
  s = s
    .replace(/([A-Z]+)/g, " $1")
    .replace(/([A-Z][a-z])/g, " $1")
    .replace("-", " ");
  return s.charAt(0).toUpperCase() + s.slice(1);
};

export let isNonEmptyList = s => {
  return typeof s == "object" && Array.isArray(s) && s.length > 0;
};

export let isObject = s => {
  return typeof s === "object";
};

export let isFunction = s => {
  return typeof s === "function";
};
