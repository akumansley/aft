<script>
import { breadcrumbStore } from '../stores.js';
import {cap} from '../util.js';
import client from '../../data/client.js';
import HLGrid from '../../ui/HLGrid.svelte';
import HLGridItem from '../../ui/HLGridItem.svelte';

let load = client.api.rpc.findMany({include: {code: true}});

breadcrumbStore.set(
	[{
		href: "/rpcs",
		text: "RPCs",
	}]
);
</script>
<HLGrid>
	{#await load}
		&nbsp;
	{:then rpcs}
		{#each rpcs as rpc}
			<HLGridItem href={"/rpc/" + rpc.id} name={rpc.name}>
			</HLGridItem>
		{/each}
		<HLGridItem href={"/rpcs/new"}>
			<div>+ Add</div>
		</HLGridItem>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLGrid>
