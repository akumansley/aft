<script>
import client from '../../data/client.js';
import {cap, getEnumsFromObj} from '../util.js';
import HLGrid from '../../ui/HLGrid.svelte';
import HLGridItem from '../../ui/HLGridItem.svelte';
import { breadcrumbStore } from '../stores.js';
let load = client.datatype.findMany({include: {validator: true, enumValues :true}});

let runtime = {}
load.then(obj => {
	for (var i = 0; i < obj.length; i++) {
		var enumValues = obj[i]["enumValues"];
		for (var j = 0; j < enumValues.length; j++) {
			runtime[enumValues[j]["id"]] = enumValues[j];
		}
	}
});

breadcrumbStore.set(
	[{
		href: "/datatypes",
		text: "Datatypes",
	}]
);
</script>
<HLGrid>
	{#await load}
		&nbsp;
	{:then datatypes}
		{#each datatypes as datatype}
			<HLGridItem href={"/datatype/" + datatype.id} name={datatype.name}>
				{#if datatype.enum == true}
				<div>Enum</div>				
				{:else}
				<div>{cap(runtime[datatype.validator.runtime]["name"])}</div>
				{/if}
			</HLGridItem>
		{/each}
		<HLGridItem href={"/datatypes/new"}>
			<div>+ Add</div>
		</HLGridItem>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLGrid>