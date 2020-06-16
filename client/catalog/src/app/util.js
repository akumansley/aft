export let cap= (s) => {
	if (!s) {
		return "";
	}
	s = s.replace(/[\w]([A-Z])/g, function(m) {
           return m[0] + " " + m[1];
       });
	return s.charAt(0).toUpperCase() + s.slice(1)
}

export let restrictToIdent= (s) => {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal.toLowerCase();
}