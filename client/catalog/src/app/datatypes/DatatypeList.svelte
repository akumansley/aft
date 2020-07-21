<script>
import client from '../../data/client.js';
import { navStore } from '../stores.js';
import {cap, getEnumsFromObj} from '../util.js';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte';
import HLRowLink from '../../ui/list/HLRowLink.svelte';
import HLBorder from '../../ui/HLBorder.svelte';
	
let load = client.api.coredatatype.findMany({include: {validator: true, enumValues :true}});

navStore.set("datatype");
let system = []
let user = []
let runtime = {}
load.then(obj => {
	for (var i = 0; i < obj.length; i++) {
		var enumValues = obj[i]["enumValues"];
		for (var j = 0; j < enumValues.length; j++) {
			runtime[enumValues[j]["id"]] = enumValues[j];
		}
		if(obj[i]["native"]) {
			system.push(obj[i]);
		} else {
			user.push(obj[i]);
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
			{#if datatype.enum == true}
			<div>Enum</div>				
			{:else}
			<div>{runtime[datatype.validator.runtime]["name"] == "starlark" ? "Code" : "Enum"}</div>
			{/if}
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
