<script>
	export let value = null;
	
	import {createEventDispatcher} from 'svelte';

	import {HLButton, HLSmallButton} from '../../ui/form/form.js';
	import {HLGrid, HLGridItem} from '../../ui/grid/grid.js';
	import {HLHeader, HLHeaderItem, HLContent, HLSectionTitle} from '../../ui/page/page.js';

	import HSpace from '../../ui/spacing/HSpace.svelte';
	import Name from '../Name.svelte';

	const dispatch = createEventDispatcher();
</script>

<HLHeader>
	<HLHeaderItem>
		<Name placeholder="Module name.." bind:value={value.name}/> 
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => {dispatch('save')}}>Save</HLButton>
	</HLHeaderItem>
</HLHeader>

<HLContent>
	<HLSectionTitle>Interfaces</HLSectionTitle>
	<HLGrid>
	{#each value.interfaces as iface}
	<HLGridItem>
		<div>{iface.name}</div>
	</HLGridItem>
	{/each}
	</HLGrid>

	<HLSectionTitle>Functions</HLSectionTitle>
	<HLGrid>
	{#each value.functions as f}
	<HLGridItem>
		<div>{f.name}</div>
	</HLGridItem>
	{/each}
	</HLGrid>

	<HLSectionTitle>Datatypes</HLSectionTitle>
	<HLGrid>
	{#each value.datatypes as dt}
	<HLGridItem>
		<div>{dt.name}</div>
	</HLGridItem>
	{/each}
	</HLGrid>
	
	<HLSectionTitle>Roles</HLSectionTitle>
	<HLGrid>
	{#each value.roles as role}
	<HLGridItem>
		<div>{role.name}</div>
	</HLGridItem>
	{/each}
	</HLGrid>
</HLContent>

