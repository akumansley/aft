<script>

	import {router} from '../router.js';
	import client from '../../data/client.js';
	import {ObjectOperation, AttributeOperation, RelationshipOperation, ConnectOperation} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';

	import EnumForm from './EnumForm.svelte';

	async function saveAndNav() {
		if (nonEmpty(value.op())) {
			let createOp = value.op().create;
			const data = await client.api.enum.create(createOp);
		}
		router.route("/datatypes");
	}

	let value = ObjectOperation({
		name: AttributeOperation(""),
		enumValues: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
			})
		),
		module: ConnectOperation(),
	});

</script>

<EnumForm on:save={saveAndNav} bind:value={value} />

