<script>
	export let value = null;
	
	import {onMount, createEventDispatcher} from 'svelte';

	import AttributeForm from './AttributeForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';

	import HLButton from '../../ui/form/HLButton.svelte';
	import HLHeader from '../../ui/page/HLHeader.svelte';
	import HLHeaderItem from '../../ui/page/HLHeaderItem.svelte';
	import HLContent from '../../ui/page/HLContent.svelte';
	import Box from '../../ui/spacing/Box.svelte';
	import Name from '../Name.svelte';
	import HLSectionTitle from '../../ui/list/HLSectionTitle.svelte';

	const dispatch = createEventDispatcher();

	function addAttribute() {
		value.attributes = value.attributes.add();
	};

	function addRelationship() {
		value.relationships = value.relationships.add();
	}

</script>

<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Model name.." bind:value={value.name}/> 
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => {dispatch('save')}}>Save</HLButton>
	</HLHeaderItem>

</HLHeader>

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
