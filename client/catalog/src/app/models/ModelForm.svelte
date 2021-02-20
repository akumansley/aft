<script>
	export let value = null;
	
	import {onMount, createEventDispatcher} from 'svelte';
	import {ConnectOperation, TypeSpecifier, ReadOnly,} from '../../api/object.js';

	import AttributeForm from './AttributeForm.svelte';
	import ImplementsForm from './ImplementsForm.svelte';
	import RelationshipForm from './RelationshipForm.svelte';
	import ReverseRelationshipForm from './ReverseRelationshipForm.svelte';
	import TargetedForm from './TargetedForm.svelte';

	import {HLButton, HLSmallButton} from '../../ui/form/form.js';
	import HSpace from '../../ui/spacing/HSpace.svelte';
	import Box from '../../ui/spacing/Box.svelte';
	import Name from '../Name.svelte';
	import ConnectSelect from '../../api/ConnectSelect.svelte';
	import {HLSectionTitle, HLHeader, HLHeaderItem, HLHeaderDetail, HLContent} from '../../ui/page/page.js';

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
	let showDetail = false;

</script>

<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Model name.." bind:value={value.name}/> 
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
	<HLSectionTitle>Interfaces</HLSectionTitle>
	<ImplementsForm bind:value={value.implements} />

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

<!--  <pre>
{JSON.stringify(value.op(), "", " ")}
</pre>
<pre>
{JSON.stringify(value, "", " ")}
</pre>
--></HLContent>

