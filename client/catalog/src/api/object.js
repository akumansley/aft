import { mergeOps } from './merge.js'

function empty(o) {
	const s = JSON.stringify(o);
	return s === "{}" || s === '[]';
};

function nonEmpty(o) {
	return !empty(o);
};

function clone(o) {
	return JSON.parse(JSON.stringify(o));
}


export function ObjectOperation(config) {
	const childOps = new Set();

	let init = null;

	const base = {
		__op: true,
		op: () => {
			const data = {};

			for (let key of childOps) {
				const childOp = base[key].op();
				if (childOp !== null && nonEmpty(childOp)) {
					data[key] = childOp;
				}
			}

			if (empty(data)) {
				return {};
			}

			let op;

			if (init) {
				op = {
					update: {
						where: {id: init.id},
						data: data,
					}
				}
			} else {
				op = {
					create: {
						data: data,
					}
				}
			}

			return op;
		},
		initialize: (i) => {
			init = JSON.parse(JSON.stringify(i));
			for (let [k, v] of Object.entries(i)) {
				if (childOps.has(k)) {
					base[k].initialize(v)
				}
			}
		},
		clone: () => {
			const newConfig = {}
			for (let [k, v] of Object.entries(config)) {
				newConfig[k] = v.clone();
			}
			return ObjectOperation(newConfig);
		}
	};
	
	for (let [k,v] of Object.entries(config)) {
		if (v.__op) {
			childOps.add(k);
		} 
		base[k] = v;
	}


	return new Proxy(base, {
		set: function(target, prop, newVal) {
			if (childOps.has(prop)) {
				const d = target[prop]
				if (d.__descriptor) {
					d.set(newVal)
					return true;
				} else if (d.__op && newVal.__op) {
					target[prop] = newVal;
					return true;
				}
			}
		},
		get: function(target, prop) {
			const v = target[prop];
			if (v && v.__descriptor) {
				return v.get();
			}
			return v;
		}
	});
}

export function RelationshipOperation(config) {
	let init = null;
	let values = [];

	function makeProxy(base, values) {
		base.__values = values;
		return new Proxy(values, {
			get: function(target, prop) {
				if (base[prop]) {
					return base[prop]
				} else {
					return target[prop]
				}
			}
		});
	}

	const base = {
		__op: true,
		op: () => {
			let ops = [];
			for (let child of base.__values) {
				const childOp = child.op();
				if (nonEmpty(childOp)) {
					ops.push(childOp);
				}
			}
			return mergeOps(ops);
		},
		initialize: (ivals) => {
			init = JSON.parse(JSON.stringify(ivals));
			for (let v of ivals) {
				const c = config.clone();
				c.initialize(v);
				values.push(c);
			}

		},
		add: () => {
			return makeProxy(base, [...base.__values, config.clone()]);
		},
		clone: () => {
			const newConfig = {}
			for (let [k, v] of Object.entries(config)) {
				newConfig[k] = v.clone();
			}
			return RelationshipOperation(newConfig);
		}
	}

	return makeProxy(base, values);
}

export function AttributeOperation(def) {
	const descriptor = {
		__op: true,
		__descriptor: true,
		init: null,
		value: def,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			descriptor.value = newVal
		},
		initialize: function(iVal) {
			descriptor.init = iVal;
			descriptor.value = iVal;
		},
		op: function() {
			if (descriptor.init !== null && descriptor.value !== descriptor.init) {
				return descriptor.value;
			} else if (descriptor.init === null) {
				return descriptor.value;
			}

			return null
		},
		clone: function() {
			return AttributeOperation(def);
		}
	}
	return descriptor;
}


export function TypeSpecifier(ifaceName) {
	const descriptor = {
		__op: true,
		__descriptor: true,
		value: ifaceName,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			return false;
		},
		initialize: function(iVal) {},
		op: function() {
			return descriptor.value;
		},
		clone: function() {
			return TypeSpecifier(ifaceName);
		}
	}
	return descriptor;
}

export function ConnectOperation() {
	const descriptor = {
		__op: true,
		__descriptor: true,
		init: null,
		value: null,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			descriptor.value = newVal;
		},
		initialize: function(iVal) {
			descriptor.init = clone(iVal);
			descriptor.value = iVal;
		},
		op: function() {
			if (descriptor.init && descriptor.value && 
				descriptor.value.id !== descriptor.init.id) {
				return {
					connect: {id: descriptor.value.id}
				}
			} else if (descriptor.init === null && descriptor.value) {
				return {
					connect: {id: descriptor.value.id}
				}
			}
			return null
		},
		clone: function() {
			return ConnectOperation();
		}

	}
	return descriptor;
}
