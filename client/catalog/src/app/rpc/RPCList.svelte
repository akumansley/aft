<script>
import { breadcrumbStore } from '../stores.js';
import {cap} from '../util.js';
import client from '../../data/client.js';
import HLGrid from '../../ui/HLGrid.svelte';
import HLGridItem from '../../ui/HLGridItem.svelte';

let load = client.rpc.findMany({include: {code: true}});

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
			<HLGridItem type={"rpc"} url={rpc.id} name={rpc.name}>
			</HLGridItem>
		{/each}
		<HLGridItem type={"rpcs"} url={"new"} name={""}>
			<div>+ Add</div>
		</HLGridItem>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLGrid>
