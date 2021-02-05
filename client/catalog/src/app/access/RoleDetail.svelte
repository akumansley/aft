<script>
	export let params = null;
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import client from '../../data/client.js';
	import {nonEmpty} from '../../lib/util.js';
	import {ObjectOperation, AttributeOperation, RelationshipOperation, SetOperation, ConnectOperation} from '../../api/object.js';
	import RoleForm from './RoleForm.svelte';
	
	navStore.set("access")
	
	let value = ObjectOperation({
		name: AttributeOperation(""),
		policies: RelationshipOperation(
			ObjectOperation({
				"interface": SetOperation(),
				readWhere: AttributeOperation("{}"),
				createWhere: AttributeOperation("{}"),
				updateWhere: AttributeOperation("{}"),
				allowRead: AttributeOperation(true),
				allowCreate: AttributeOperation(true),
				allowUpdate: AttributeOperation(true),
			}),
			),
		executableFunctions: RelationshipOperation(ConnectOperation()),
		module: SetOperation(),
	});

	let load = client.api.role.findOne({
		where: {id: params.id}, 
		include: {
			executableFunctions: true,
			policies: {
				include: { 
					"interface": true 
				},
			},
			module: true,
		},
	}).then((data) => {
		value.initialize(data);
		value = value;
	});


	async function save() {
		const op = value.op();
		if (nonEmpty(op)) {
			await client.api.role.update(op.update);
		}
		router.route('/roles');
	}


</script>

{#await load then loaded}
<RoleForm on:save={save} bind:value={value} />
{/await}
