<script>
	import client from '../../data/client.js';
	import { navStore } from '../stores.js';
	navStore.set("datatype");
	
	import {cap } from '../../lib/util.js';

	import HLGrid from '../../ui/grid/HLGrid.svelte';
	import HLGridItem from '../../ui/grid/HLGridItem.svelte';
	import HLGridNew from '../../ui/grid/HLGridNew.svelte';
	import HLSectionTitle from '../../ui/page/HLSectionTitle.svelte';
	import HLBorder from '../../ui/page/HLBorder.svelte';
	import HLContent from '../../ui/page/HLContent.svelte';
	
	let load = client.api.datatype.findMany({});

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

	function urlFor(dt) {
		switch (dt.type) {
			case "enum":
			return "/enum/" + dt.id;
			case "coreDatatype":
			return "/core/" + dt.id;
		}
	}
</script>

<style>
	.v-space {
		height: var(--box-margin);
	}
</style>

{#await load then load}
<HLGrid>
	<HLGridNew href={"/enums/new"}>Add Enum</HLGridNew>
</HLGrid>
<HLBorder/>
<HLContent>
<HLSectionTitle>Datatypes</HLSectionTitle>
<HLGrid>
	{#each user as datatype}
	<HLGridItem href={"/datatype/" + datatype.id} name={datatype.name}>
	</HLGridItem>
	{/each}
</HLGrid>
<div class="v-space"></div>
<HLBorder/>
<div class="v-space"></div>
<HLSectionTitle>System</HLSectionTitle>
<HLGrid>
	{#each system as datatype}
	<HLGridItem name={datatype.name} href={urlFor(datatype)}>
	</HLGridItem>
	{/each}
</HLGrid>
</HLContent>
{/await}

