<script>
	export let params = null;

	import { router } from '../router.js';
	import client from '../../data/client.js';
	import {ObjectOperation, AttributeOperation, RelationshipOperation, ConnectOperation} from '../../api/object.js';
	import EnumForm from './EnumForm.svelte';
	import {nonEmpty} from '../../lib/util.js';

	async function saveAndNav() {
		if (nonEmpty(value.op())) {
			let updateOp = value.op().update;
			const data = await client.api.enum.update(updateOp);
		}
		router.route("/datatypes");
	}

	let value = ObjectOperation({
		name: AttributeOperation(""),
		enumValues: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
			})),
		module: ConnectOperation(),
	})

	let load = client.api.enum.findOne({
		where: {id: params.id},
		include: {enumValues: true, module: true},
	}).then((e) => {
		value.initialize(e);
		value = value;
	});
</script>

<EnumForm on:save={saveAndNav} bind:value={value} />

