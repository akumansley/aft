<script>
	import client from '../../data/client.js';
	import {cap} from '../../lib/util.js';
	import {navStore} from '../stores.js';

	import {HLGrid, HLGridItem, HLGridNew, HLSecondary} from '../../ui/grid/grid.js';
	import {HLBorder, HLContent, HLSectionTitle, HLCallout, HLPad} from '../../ui/page/page.js';

	navStore.set("schema");

	let appModule = null;
	let nativeModules = [];
	let load = client.api.module.findMany({
		where:{OR:[
			{interfaces:{some:{}}},
			{goPackage: ""},
			]},
			include: {interfaces: true}})
	.then(result => {
		for (let mod of result) {
			if (mod.goPackage === "") {
				appModule = mod;
			} else {
				nativeModules.push(mod);
			}
		}
	});

	function urlFor(iface) {
		if (iface.type === "concreteInterface") {
			return "/interface/" + iface.id
		} else {
			return "/model/" + iface.id
		}
	}
</script>

{#await load then load}
<HLGrid>
	<HLGridNew href={"/models/new"}>Add Model</HLGridNew>
	<HLGridNew href={"/interfaces/new"}>Add Interface</HLGridNew>
</HLGrid>
<HLBorder />

<HLContent>
	{#if appModule.interfaces.length}
	<HLPad>
		<HLSectionTitle>Models</HLSectionTitle>
		<HLGrid>
			{#each appModule.interfaces as iface}
			<HLGridItem href={urlFor(iface)}>
				<div>{cap(iface.name)}</div>
				<HLSecondary>{appModule.name}</HLSecondary>
			</HLGridItem>
			{/each}
		</HLGrid>
	</HLPad>
	{:else}
	<HLCallout>
		<div>No application models yet. Add some to get started!</div>
	</HLCallout>
	{/if}

	<HLBorder spaceBottom={true} />
	<HLSectionTitle>System</HLSectionTitle>
	<HLGrid>
		{#each nativeModules as mod}
		{#each mod.interfaces as iface}
		<HLGridItem href={urlFor(iface)}>
			<div>{cap(iface.name)}</div>
			<HLSecondary>{mod.name}</HLSecondary>
		</HLGridItem>
		{/each}
		{/each}
	</HLGrid>
</HLContent>
{/await}
