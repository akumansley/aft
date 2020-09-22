<script>
	export let value = null;
	
	import {onMount, createEventDispatcher} from 'svelte';
	import {ConnectOperation, TypeSpecifier, ReadOnly,} from '../../api/object.js';

	import AttributeForm from './AttributeForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';
	import ReverseRelationshipForm from './ReverseRelationshipForm.svelte';
	import TargetedForm from './TargetedForm.svelte';

	import {HLButton, HLSmallButton} from '../../ui/form/form.js';
	import HLHeader from '../../ui/page/HLHeader.svelte';
	import HLHeaderItem from '../../ui/page/HLHeaderItem.svelte';
	import HSpace from '../../ui/spacing/HSpace.svelte';
	import HLContent from '../../ui/page/HLContent.svelte';
	import Box from '../../ui/spacing/Box.svelte';
	import Name from '../Name.svelte';
	import HLSectionTitle from '../../ui/list/HLSectionTitle.svelte';

	const dispatch = createEventDispatcher();

	function addAttribute() {
		value.attributes = value.attributes.add();
	}

	function addRelationship() {
		value.relationships = value.relationships.add({
			type: "concreteRelationship",
		});
	}

	function addReverseRelationship(referencing) {
		value.relationships = value.relationships.add({
			type: "reverseRelationship",
			referencing: referencing,
		});
	}

	function notReversed(tRel) {
		for (let rel of value.relationships) {
			if (rel.type === "reverseRelationship" && rel.referencing.id === tRel.id) {
				return false;
			}
		}
		return true;
	}

	function anyNotReversed() {
		for (let tRel of value.targeted) {
			if (notReversed(tRel)) {
				return true;
			}
		}
		return false;
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
	{#if rel.type === "concreteRelationship"}
	<RelationshipForm bind:value={rel} />
	{/if}

	{#if rel.type === "reverseRelationship"}
	<ReverseRelationshipForm bind:value={rel} />
	{/if}
	{/each}
	<Box>
		<HLButton on:click={addRelationship}>+add</HLButton>
	</Box>

	{#if anyNotReversed(value.targeted, value.relationships)}
	<HLSectionTitle>Referenced by</HLSectionTitle>
	{#each value.targeted as tRel}
	{#if notReversed(tRel)}
	<TargetedForm on:click={() => addReverseRelationship(tRel)} value={tRel} />
		{/if}
		{/each}
		{/if}

 <pre>
	{JSON.stringify(value.op(), "", " ")}
</pre>
 <pre>
	{JSON.stringify(value, "", " ")}
</pre>
</HLContent>

