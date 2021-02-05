<script>
	import client from '../../data/client.js';
	import { navStore } from '../stores.js';
	navStore.set("datatype");
	
	import {cap } from '../../lib/util.js';

	import {HLGrid, HLGridItem, HLGridNew, HLSecondary} from '../../ui/grid/grid.js';
	import {HLSectionTitle, HLBorder, HLContent, HLCallout, HLPad} from '../../ui/page/page.js';

	navStore.set("datatype");
	
	let appModule = null;
	let nativeModules = [];
	let load = client.api.module.findMany({
		where:{ OR:[
			{datatypes:{some:{}}},
			{goPackage: ""},
			]}, 
			include: {datatypes: true}})
	.then(result => {
		for (let mod of result) {
			if (mod.goPackage === "") {
				appModule = mod;
			} else {
				nativeModules.push(mod);
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

{#await load then load}
<HLGrid>
	<HLGridNew href={"/enums/new"}>Add Enum</HLGridNew>
</HLGrid>
<HLBorder/>

<HLContent>
	{#if appModule.datatypes.length}
	<HLPad>
		<HLSectionTitle>Datatypes</HLSectionTitle>
		<HLGrid>
			{#each appModule.datatypes as dt}
			<HLGridItem href={urlFor(dt)}>
				<div>{cap(dt.name)}</div>
				<HLSecondary>{appModule.name}</HLSecondary>
			</HLGridItem>
			{/each}
		</HLGrid>
	</HLPad>
	{:else}
	<HLCallout>
		<div>No application datatypes yet</div>
	</HLCallout>
	{/if}

	<HLBorder spaceBottom={true} />
	<HLSectionTitle>System</HLSectionTitle>

	<HLGrid>
		{#each nativeModules as mod}
		{#each mod.datatypes as dt}
		<HLGridItem href={urlFor(dt)}>
			<div>{cap(dt.name)}</div>
			<HLSecondary>{mod.name}</HLSecondary>
		</HLGridItem>
		{/each}
		{/each}
	</HLGrid>
</HLContent>
{/await}

