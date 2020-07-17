export let checkSave= (func) => {
	return (e) => {
		if ((e.metaKey || e.ctrlKey) && e.keyCode == 83) { /*ctrl+s or command+s*/
			e.preventDefault();
			func();
			return false;
		}
	}
}