export function mergeOps(values) { 
	let op = {};
	for (let value of values) {
		if (!value) {
			continue;
		}
		if (value.create) {
			if (op.create) {
				op.create = [...op.create, value.create.data];
			} else {
				op.create = [value.create.data];
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
		} else if (value.connect) {
			if (op.connect) {
				op.connect = [...op.connect, value.connect];
			} else {
				op.connect = [value.connect];
			}
		}
	}
	return op;
}