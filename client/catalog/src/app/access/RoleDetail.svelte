<script>
	export let params = null;
	import {navStore} from '../stores.js';
	import {router} from '../router.js';

	import client from '../../data/client.js';

	import HLButton from '../../ui/form/HLButton.svelte';
	import {HLHeader, HLHeaderItem, HLContent} from '../../ui/page/page.js';
	import {Box} from '../../ui/spacing/spacing.js';

	import Name from '../Name.svelte';
	import RolesPicker from './RolesPicker.svelte';
	import PolicyForm from './PolicyForm.svelte';

	import HLSectionTitle from '../../ui/page/HLSectionTitle.svelte';
	import {nonEmpty} from '../../lib/util.js';
	import {ObjectOperation, AttributeOperation, RelationshipOperation, SetOperation} from '../../api/object.js';


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
	});

	let load = client.api.role.findOne({
		where: {id: params.id}, 
		include: {
			policies: {
				include: { 
					"interface": true 
				},
			},
		},
	}).then((data) => {
		value.initialize(data);
		value = value;
	});

	function addPolicy() {
		value.policies = value.policies.add();
	};

	async function save() {
		const op = value.op();
		if (nonEmpty(op)) {
			await client.api.role.update(op.update);
		}
		router.route('/roles');
	}


</script>

{#await load then loaded}

<HLHeader>
	<HLHeaderItem>		
		<Name placeholder="Role name.." bind:value={value.name} />
	</HLHeaderItem>		
	<HLHeaderItem>	
		<HLButton on:click={save}>Save</HLButton>
	</HLHeaderItem>
</HLHeader>
	
<HLContent>
	<HLSectionTitle>Grants</HLSectionTitle>

	{#each value.policies as policy}
	<PolicyForm bind:value={policy} />
	{/each}

	<Box>
		<HLButton on:click={addPolicy}>+ add</HLButton>
	</Box>
</HLContent>

{/await}
