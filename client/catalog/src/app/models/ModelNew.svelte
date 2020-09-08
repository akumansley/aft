<script>
	import { onMount } from 'svelte';
	import { navStore } from '../stores.js';
	import { router } from '../router.js';

	import client from '../../data/client.js';
	import AttributeForm from './AttributeForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';

	import {Create} from '../../api/api.js';
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

	let model = {
		name: "",
		attributes: {create: []},
		relationships: {create: []},
	};

	function addAttribute() {
		model.attributes.create = [...model.attributes.create, {}];
	}

	function addRelationship() {
		model.relationships.create = [...model.relationships.create, {}];
	}

	async function saveAndNav() {
		const data = await client.api.model.create({data: model});
		router.route("/models/" + data.id);
	}


</script>

<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Model name.." bind:value={model.name} />
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={saveAndNav}>Save</HLButton>
	</HLHeaderItem>
</HLHeader>

<HLContent>
	<Create>
		<HLSectionTitle>Attributes</HLSectionTitle>
		{#each model.attributes.create as attr}
		<AttributeForm bind:value={attr}/>
		{/each}
		<Box>
			<HLButton on:click={addAttribute}>+add</HLButton>
		</Box>

		<HLSectionTitle>Relationships</HLSectionTitle>
		{#each model.relationships.create as rel}
		<RelationshipForm bind:value={rel}/>
		{/each}
		<Box>
			<HLButton on:click={addRelationship}>+add</HLButton>
		</Box>
	</Create>
</HLContent>
