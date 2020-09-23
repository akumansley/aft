<script>
import {cap} from '../../lib/util.js';
import { navStore } from '../stores.js';
import client from '../../data/client.js';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte'
import HLRowLink from '../../ui/list/HLRowLink.svelte';
import HLBorder from '../../ui/page/HLBorder.svelte'
import HLListTitle from '../../ui/list/HLListTitle.svelte';
import HLSectionTitle from '../../ui/page/HLSectionTitle.svelte';

let rpcs = client.api.rpc.findMany({include: {code: true}});
navStore.set("rpc");
</script>

<style>
	.v-space {
		height: var(--box-margin);
	}
</style>

{#await rpcs then rpcs}
<HLListTitle>RPCs</HLListTitle>
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
<HLSectionTitle>System</HLSectionTitle>
	<HLGrid>
{#each rpcs as rpc}
	{#if rpc.native == true}
	<HLGridItem name={rpc.name} href={"/rpc/" + rpc.id}/ >
	</HLGridItem>
	{/if}
{/each}
</HLGrid>
{/await}
