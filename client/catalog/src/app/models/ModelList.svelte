<script>
import client from '../../data/client.js';
import { breadcrumbStore } from '../stores.js';
import HLGrid from '../../ui/HLGrid.svelte';
import HLGridItem from '../../ui/HLGridItem.svelte';
let load = client.model.findMany({include: {attributes: true}});

breadcrumbStore.set(
	[{
		href: "/models",
		text: "Models",
	}]
);
</script>
<HLGrid>
	{#await load}
		&nbsp;
	{:then models}
		{#each models as model}
			<HLGridItem href={"/model/" + model.id} name={model.name}>
				{#each model.attributes as attr}
					<div>{attr.name}</div>
				{/each}
			</HLGridItem>
		{/each}
		<HLGridItem href={"/models/new"}>
			<div>+ Add</div>
		</HLGridItem>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLGrid>
