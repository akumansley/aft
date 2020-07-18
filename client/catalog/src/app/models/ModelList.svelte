<script>
import client from '../../data/client.js';
import { cap } from '../util.js';
import { navStore } from '../stores.js';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte';
import HLRowLink from '../../ui/list/HLRowLink.svelte';
import HLBorder from '../../ui/HLBorder.svelte';

let load = client.api.model.findMany({include: {attributes: true}});
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
navStore.set("model");
</script>

<style>
	.v-space {
		height: var(--box-margin);
	}
</style>

{#await load then load}
<h1>Models</h1>
<HLGrid>
	<HLGridNew href={"/models/new"} />
	{#each user as model}
		<HLGridItem href={"/model/" + model.id} name={model.name}>
			{#each model.attributes as attr}
				<div>{cap(attr.name)}</div>
			{/each}
		</HLGridItem>
	{/each}
</HLGrid>
<div class="v-space"></div>
<HLBorder />
<div class="v-space"></div>
<h2>System</h2>
{#each system as model}
	<HLRowLink href={"/model/" + model.id}>
		{cap(model.name)}
	</HLRowLink>
{/each}
{/await}