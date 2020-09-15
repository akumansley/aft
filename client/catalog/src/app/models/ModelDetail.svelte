<script>
	export let params = null;
	import { onMount } from 'svelte';
	import { navStore } from '../stores.js';
	import {router} from '../router.js';

	import client from '../../data/client.js';
	import AttributeForm from './AttributeForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';

	import {Many, mergeOps} from '../../api/api.js';

	import {ObjectOperation, RelationshipOperation, AttributeOperation, ConnectOperation} from '../../api/object.js';

	import { clone, nonEmpty } from '../../lib/util.js';
	import HLRowButton from '../../ui/list/HLRowButton.svelte';
	import HLButton from '../../ui/form/HLButton.svelte';
	import HLRow from '../../ui/list/HLRow.svelte';
	import HLHeader from '../../ui/page/HLHeader.svelte';
	import HLHeaderItem from '../../ui/page/HLHeaderItem.svelte';
	import HLContent from '../../ui/page/HLContent.svelte';
	import Box from '../../ui/spacing/Box.svelte';
	import Name from '../Name.svelte';
	import HLSectionTitle from '../../ui/list/HLSectionTitle.svelte';

	navStore.set("schema");


	let init;
	let model = ObjectOperation({
		name: AttributeOperation(""),
		relationships: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
				multi: AttributeOperation(false),
				target: ConnectOperation(),
			})),
		attributes: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
				datatype: ConnectOperation(),
			}))
	});

	let load = client.api.model.findOne({
		where: {id: params.id},
		include: {
			attributes: {
				include: {datatype: true},
			},
			relationships: {
				case: {
					concreteRelationship: {
						include: {target: true},
					}
				}
			}
		},
	}).then(m => { 
		try {
			model.initialize(m);
			model = model;
		} catch (e) {
			console.log(e);
		}
	});


	function addAttribute() {
		model.attributes.add();
		model = model;
	};

	function addRelationship() {
		model.relationships.add();
		model = model;
	}

	async function saveAndNav() {
		await save();
		router.route("/schema");
	}

	async function save() {
		return client.api.model.update(model.op());
	}

</script>
<style>
	pre {
		font-size:13px;
	}
</style>

{#await load then _}
<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Model name.." bind:value={model.name}/> 
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={saveAndNav}>Save</HLButton>
	</HLHeaderItem>

</HLHeader>

<HLContent>
	<HLSectionTitle>Attributes</HLSectionTitle>
	
	{#each model.attributes as attr}
	<AttributeForm bind:value={attr} />
	{/each}

	<Box>
		<HLButton on:click={addAttribute}>+add</HLButton>
	</Box>
 	<HLSectionTitle>Relationships</HLSectionTitle>
	{#each model.relationships as rel}
	<RelationshipForm bind:value={rel} />
	{/each}
	<Box>
		<HLButton on:click={addRelationship}>+add</HLButton>
	</Box>
	
<pre style="float: left;">Op:
{JSON.stringify(model.op(), null, 2)}
</pre>

<pre style="float: left; margin-left: 2em;">Value:
{JSON.stringify(model, null, 2)}
</pre>
</HLContent>

{/await}