<script>
	export let value = null;
	
	import {onMount, createEventDispatcher} from 'svelte';

	import AttributeForm from './AttributeForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';

	import {HLButton, HLSmallButton} from '../../ui/form/form.js';
	import {HLHeader, HLHeaderItem, HLHeaderDetail, HLContent, HLSectionTitle} from '../../ui/page/page.js';
	import HSpace from '../../ui/spacing/HSpace.svelte';
	import Box from '../../ui/spacing/Box.svelte';
	import Name from '../Name.svelte';
	import ConnectSelect from '../../api/ConnectSelect.svelte';

	const dispatch = createEventDispatcher();

	function addAttribute() {
		value.attributes = value.attributes.add();
	}

	function addRelationship() {
		value.relationships = value.relationships.add({
			type: "concreteRelationship",
		});
	}

	let showDetail = false;

</script>

<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Interface name.." bind:value={value.name}/> 
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => {dispatch('save')}}>Save</HLButton>
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => showDetail = !showDetail}>More</HLButton>
	</HLHeaderItem>

</HLHeader>

<HLHeaderDetail show={showDetail}>
	<HLHeaderItem>
		Module: <HSpace/> <ConnectSelect pickDefault={(m) => m.goPackage === ""} bind:value={value.module} iface={"module"} />
	</HLHeaderItem>
</HLHeaderDetail>

<HLContent>
	<HLSectionTitle>Attributes</HLSectionTitle>
	
	{#each value.attributes as attr}
	<AttributeForm bind:value={attr} />
	{/each}

	<Box>
		<HLButton on:click={addAttribute}>+add</HLButton>
	</Box>
	<HLSectionTitle>Relationships</HLSectionTitle>
	{#each value.relationships as rel}
	<RelationshipForm bind:value={rel} />
	{/each}
	<Box>
		<HLButton on:click={addRelationship}>+add</HLButton>
	</Box>
</HLContent>

