<script>
import {cap} from '../util.js';
import { navStore } from '../stores.js';
import aft from '../../data/aft.js';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte'
import HLRowLink from '../../ui/list/HLRowLink.svelte';
import HLBorder from '../../ui/HLBorder.svelte'

let rpcs = aft.api.rpc.findMany({include: {code: true}});
navStore.set("rpc");
</script>

<style>
	.v-space {
		height: var(--box-margin);
	}
</style>

{#await rpcs then rpcs}
<h1>Functions</h1>
<HLGrid>
	<HLGridNew href={"/rpcs/new"}/>
	{#each rpcs as rpc}
		{#if rpc.native == false}
		<HLGridItem href={"/rpc/" + rpc.id} name={rpc.name}>
		</HLGridItem>
		{/if}
	{/each}
</HLGrid>
<div class="v-space"></div>
<HLBorder/>
<div class="v-space"></div>
<h2>System</h2>
{#each rpcs as rpc}
	{#if rpc.native == true}
	<HLRowLink href={"/rpc/" + rpc.id}>
		{cap(rpc.name)}
	</HLRowLink>
	{/if}
{/each}
{/await}
