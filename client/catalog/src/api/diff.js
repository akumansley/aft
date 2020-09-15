export async function diffObject(init, value, type, client) {
	let schema = await client.schema[type];

	if (init) {
		op = {
			update: {
				where: {
					id: init.id,
				},
				data: {
					...diffAttributesUpdate(init, value, schema, client),
					...diffRelationshipsUpdate(init, value, schema, client)
				}
			}
		}
	} else {
		op = {
			create: {
				...diffAttributesCreate(init, value, schema, client),
				...diffRelationshipsCreate(init, value, schema, client)
			}
		}
	}

	return op;
}

async function diffAttributesUpdate(init, value, schema, client) {
	const op = {};
	for (let attr of schema.attributes) {
		const key = attr.name;
		const left = init[key];
		const right = value[key];

		if (left !== right) {
			op[key] = right;
		}
	}
	return op;
}


async function diffRelationshipsUpdate(init, value, schema, client) {
	const op = {};
	for (let rel of schema.relationships) {
		const key = rel.name;
		const left = init[key];
		const right = value[key];

		// left and right are sets
		// hash the left ids and compare
		if (rel.multi) {
			const leftIds = new Set();

			const created = new Set();
			const preserved = new Set();
			const connected = new Set();

			for (let lo of left) {
				leftIds.add(lo.id);
			}

			for (let ro of right) {
				// we use the presence of an id field
				// to distinguish between connecting to
				// an already saved object
				// and creating a new one on the fly
				if (!ro.id) {
					created.add(ro);
				} else if (leftIds.has(ro.id)) {
					leftIds.remove(ro.id)
					preserved.add(ro);
				} else {
					connected.add(ro);
				}
			}

			// those that remain..
			// is this too janky? we need tombstones to do it right..
			const disconnected = leftIds;
		}
		// left and right are single objects


	}
	return op;
}