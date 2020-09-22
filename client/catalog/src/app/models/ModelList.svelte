<script>
import client from '../../data/client.js';
import { cap } from '../../lib/util.js';
import { navStore } from '../stores.js';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte';
import HLRowLink from '../../ui/list/HLRowLink.svelte';
import HLBorder from '../../ui/page/HLBorder.svelte';
import HLListTitle from '../../ui/list/HLListTitle.svelte';
import HLSectionTitle from '../../ui/list/HLSectionTitle.svelte';

let load = client.api.interface.findMany({});
var system = [];
var user = [];
load.then(obj => {
	for(var i = 0; i < obj.length; i++) {
		if(obj[i].system === true) {
			system.push(obj[i]);
		} else {
			user.push(obj[i]);
		}
	}
});
navStore.set("schema");
</script>

<style>
	.v-space {
		height: var(--box-margin);
	}
</style>

{#await load then load}
<HLListTitle>Schema</HLListTitle>
<HLGrid>
	<HLGridNew href={"/models/new"}>Add Model</HLGridNew>
		<HLGridNew href={"/interfaces/new"}>Add Interface</HLGridNew>
	{#each user as iface}
	{#if iface.type === "model"}
		<HLGridItem href={"/model/" + iface.id} name={iface.name}></HLGridItem>
	{/if}
	{#if iface.type === "concreteInterface"}
		<HLGridItem href={"/interface/" + iface.id} name={iface.name}></HLGridItem>
	{/if}
	{/each}
</HLGrid>
<div class="v-space"></div>
<HLBorder />
<div class="v-space"></div>
<HLSectionTitle>System</HLSectionTitle>
<HLGrid>
{#each system as model}
	<HLGridItem name={model.name} href={"/model/" + model.id}>
	</HLGridItem>
{/each}
</HLGrid>
{/await}
