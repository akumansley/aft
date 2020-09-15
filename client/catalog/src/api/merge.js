export function mergeOps(values) { 
	let op = {};
	for (let value of values) {
		if (!value) {
			continue;
		}
		if (value.create) {
			if (op.create) {
				op.create = [...op.create, value.create];
			} else {
				op.create = [value.create];
			}
		} else if (value.update) {
			if (op.update) {
				op.update = [...op.update, value.update];
			} else {
				op.update = [value.update];
			}
		} else if (value.delete) {
			if (op.delete) {
				op.delete = [...op.delete, value.delete];
			} else {
				op.delete = [value.delete];
			}
		}
	}
	return op;
}