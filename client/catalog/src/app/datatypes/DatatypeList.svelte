<script>
import aft from '../../data/aft.js';
import { navStore } from '../stores.js';
import {cap } from '../util.js';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte';
import HLRowLink from '../../ui/list/HLRowLink.svelte';
import HLBorder from '../../ui/HLBorder.svelte';
	
let load = aft.api.datatype.findMany({});

navStore.set("datatype");
let system = []
let user = []
let runtime = {}
load.then(dts => {
	for (let dt of dts) {
		if(dt.system) {
			system.push(dt);
		} else {
			user.push(dt);
		}
	}
});
</script>

<style>
	.v-space {
		height: var(--box-margin);
	}
</style>

{#await load then load}
<h1>Datatypes</h1>
	<HLGrid>
		<HLGridNew href={"/datatypes/new"}/>
		{#each user as datatype}
		<HLGridItem href={"/datatype/" + datatype.id} name={datatype.name}>
		</HLGridItem>
		{/each}
	</HLGrid>
<div class="v-space"></div>
<HLBorder/>
<div class="v-space"></div>
<h2>System</h2>
{#each system as datatype}
	<HLRowLink href={"/datatype/" + datatype.id}>
		{cap(datatype.name)}
	</HLRowLink>
{/each}
{/await}
