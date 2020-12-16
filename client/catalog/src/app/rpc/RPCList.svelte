<script>
	import {cap} from '../../lib/util.js';
	import { navStore } from '../stores.js';
	import client from '../../data/client.js';

	import HLGrid from '../../ui/grid/HLGrid.svelte';
	import HLGridItem from '../../ui/grid/HLGridItem.svelte';
	import HLGridNew from '../../ui/grid/HLGridNew.svelte'
	import HLBorder from '../../ui/page/HLBorder.svelte'
	import HLContent from '../../ui/page/HLContent.svelte'
	import HLSectionTitle from '../../ui/page/HLSectionTitle.svelte';

	const RPC = "8decedba-555b-47ca-a232-68100fbbf756";
	let rpcs = client.api.rpc.findMany({
		where: {},
		include: {
			function: true,
		}
	});

	navStore.set("rpc");
</script>

{#await rpcs then rpcs}
<HLGrid>
	<HLGridNew href={"/rpcs/new"}>Add RPC</HLGridNew>
</HLGrid>
<HLBorder/>
<HLContent>
	<HLSectionTitle>RPCs</HLSectionTitle>
	<HLGrid>
		{#each rpcs as rpc}
		<HLGridItem href={"/rpc/" + rpc.id} name={rpc.function.name}>
		</HLGridItem>
		{/each}
	</HLGrid>
</HLContent>
{/await}