<script>
	export let params = null;
	import { onMount } from 'svelte';
	import { navStore } from '../stores.js';
	import {router} from '../router.js';

	import client from '../../data/client.js';
	import AttributeForm from './AttributeForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';

	import {Update, Create, Many} from '../../api/api.js';
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


	let model;
	let load = client.api.model.findOne({
		where: {id: params.id},
		include: {
			attributes: {
				include: {datatype: true},
			},
			relationships: true,
		},
	}).then(m => { model = m; });

	let op = {
		attributes: {
			update: [],
			create: [],
		},
	};


	let addAttribute;

	function addRelationship() {
		model.relationships = [...model.relationships, {}];
	}

	async function saveAndNav() {
		await save();
		router.route("/models");
	}

	async function save() {}

</script>

{#await load then _}
<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Model name.." bind:value={model.name} on:click={saveAndNav}/> 
	</HLHeaderItem>
</HLHeader>

<HLContent>
	<HLSectionTitle>Attributes</HLSectionTitle>
	
	<Many component={AttributeForm} init={model.attributes} bind:op={op.attributes} bind:add={addAttribute}/>

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
{JSON.stringify(op, null, 2)}
	</pre>
</HLContent>

{/await}